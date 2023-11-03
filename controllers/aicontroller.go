package controllers

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/sashabaranov/go-openai"
)

type AiResponse struct {
	Status string `json:"status"`
	Data   string `json:"data"`
}

type AiUsecase interface {
	Rekomendasiproperti(pesanPengguna, kunciOpenAI string) (string, error)
}

type aiUsecase struct{}

func NewAiUsecase() AiUsecase {
	return &aiUsecase{}
}

func (uc *aiUsecase) Rekomendasiproperti(pesanPengguna, kunciOpenAI string) (string, error) {
	ctx := context.Background()
	client := openai.NewClient(kunciOpenAI)
	model := openai.GPT3Dot5Turbo
	pesan := []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleSystem,
			Content: "Halo, perkenalkan saya sewakeun sistem untuk menyewa rumah lebih aman dan efisien",
		},
		{
			Role:    openai.ChatMessageRoleUser,
			Content: pesanPengguna,
		},
	}

	resp, err := uc.getDokumenLengkapDariPesan(ctx, client, pesan, model)
	if err != nil {
		return "", err
	}
	jawaban := resp.Choices[0].Message.Content
	return jawaban, nil
}

func (uc *aiUsecase) getDokumenLengkapDariPesan(
	ctx context.Context,
	client *openai.Client,
	pesan []openai.ChatCompletionMessage,
	model string,
) (openai.ChatCompletionResponse, error) {
	if model == "" {
		model = openai.GPT3Dot5Turbo
	}

	resp, err := client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model:    model,
			Messages: pesan,
		},
	)
	return resp, err
}

func RekomendasiPropertiChatBot(c echo.Context, aiUsecase AiUsecase) error {
	tokenString := c.Request().Header.Get("Authorization")
	if tokenString == "" {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{"error": true, "message": "Token otorisasi tidak ditemukan"})
	}

	parts := strings.SplitN(tokenString, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{"error": true, "message": "Format token tidak valid"})
	}

	tokenString = parts[1]

	var dataPermintaan map[string]interface{}
	err := c.Bind(&dataPermintaan)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": true, "message": "Format JSON tidak valid"})
	}

	pesanPengguna, ok := dataPermintaan["pesan"].(string)
	if !ok || pesanPengguna == "" {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": true, "message": "Pesan tidak ditemukan atau tidak valid"})
	}

	pesanPengguna = fmt.Sprintf("Tips menyewa properti pilihan anak muda: %s", pesanPengguna)

	jawaban, err := aiUsecase.Rekomendasiproperti(pesanPengguna, os.Getenv("ALTA_MINI_PROJECT"))
	if err != nil {
		pesanKesalahan := "Gagal mendapatkan rekomendasi properti"
		if strings.Contains(err.Error(), "pembatasan API") {
			pesanKesalahan = "Pembatasan limit API. Coba lagi nanti"
		}
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": true, "message": pesanKesalahan})
	}

	dataRespons := AiResponse{
		Status: "sukses",
		Data:   jawaban,
	}

	return c.JSON(http.StatusOK, dataRespons)
}
