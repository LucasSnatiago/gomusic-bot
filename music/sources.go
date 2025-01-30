package music

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"

	"github.com/lrstanley/go-ytdlp"
)

func SetupYtdlp() {
	// If yt-dlp isn't installed yet, download and cache it for further use.
	ytdlp.MustInstall(context.TODO(), nil)
}

func YtGetBytes(url string) []byte {
	buf, err := ytGetBytes(url)
	if err != nil {
		fmt.Printf("Error getting data from: %s\n", url)
	}

	dcaBytes, err := convertToDCA(buf)
	if err != nil {
		fmt.Printf("Error converting data to DCA: %s\n", err.Error())
	}

	return dcaBytes
}

func ytGetBytes(url string) ([]byte, error) {
	cmd := exec.Command("yt-dlp", "-f", "bestaudio",
		"--extract-audio",        // Extrai apenas o áudio
		"--audio-format", "opus", // Define o formato de saída como Opus
		"--audio-quality", "0", // Qualidade máxima (0 é o melhor)
		"--sponsorblock-remove", "all", "-o", "-", url)

	var buf bytes.Buffer
	cmd.Stdout = &buf

	// Executa o comando
	err := cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("erro ao executar yt-dlp: %w", err)
	}

	// Retorna os bytes capturados
	return buf.Bytes(), nil
}

func convertToDCA(input []byte) ([]byte, error) {
	// Cria um comando ffmpeg para converter o áudio para DCA
	cmd := exec.Command(
		"ffmpeg",
		"-i", "pipe:0", // Entrada a partir de stdin
		"-f", "s16le", // Formato de saída: PCM 16-bit little-endian
		"-ar", "48000", // Taxa de amostragem: 48kHz
		"-ac", "2", // Canais de áudio: estéreo
		"pipe:1", // Saída para stdout
	)

	// Configura a entrada e a saída do comando
	cmd.Stdin = bytes.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out

	// Executa o comando
	err := cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("erro ao converter áudio para DCA: %w", err)
	}

	return out.Bytes(), nil
}
