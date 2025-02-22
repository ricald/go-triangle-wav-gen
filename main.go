package main

import (
	"encoding/binary"
	"fmt"
	"math"
	"os"
)

const (
	sampleRate = 44100
	frequency  = 440.0 * 7
	// frequency = 440.0
	amplitude = 0.5
	duration  = 5 // 秒
	channels  = 1
	bitDepth  = 16 // ビット/サンプル
)

// WAVヘッダ構造体
type wavHeader struct {
	ChunkID       [4]byte // "RIFF"
	ChunkSize     uint32  // ファイルサイズ - 8
	Format        [4]byte // "WAVE"
	Subchunk1ID   [4]byte // "fmt "
	Subchunk1Size uint32  // 16
	AudioFormat   uint16  // 1 (PCM)
	NumChannels   uint16  // チャンネル数
	SampleRate    uint32  // サンプルレート
	ByteRate      uint32  // サンプルレート * チャンネル数 * ビット/サンプル / 8
	BlockAlign    uint16  // チャンネル数 * ビット/サンプル / 8
	BitsPerSample uint16  // ビット/サンプル
	Subchunk2ID   [4]byte // "data"
	Subchunk2Size uint32  // データサイズ
}

func main() {
	// WAVファイルを作成
	file, err := os.Create("triangle.wav")
	if err != nil {
		fmt.Println("ファイルを作成できませんでした:", err)
		return
	}
	defer file.Close()

	// WAVヘッダを書き込み
	writeWavHeader(file)

	// 三角波を生成して書き込み
	for i := 0; i < sampleRate*duration; i++ {
		time := float64(i) / sampleRate
		// 鳴動と停止を切り替える
		if math.Mod(time, 1.0) < 0.5 {
			wave := amplitude * (2*math.Abs(math.Mod(2*frequency*time, 2)-1) - 1)
			// wave := amplitude * math.Sin(2*math.Pi*frequency*time)
			// var wave float64
			// if math.Mod(2*frequency*time, 2) < 1 {
			// 	wave = amplitude
			// } else {
			// 	wave = -amplitude
			// }
			writeSample(file, float32(wave))
		} else {
			writeSample(file, 0.0) // 停止時は0を出力
		}
	}

	fmt.Println("triangle.wav を作成しました。")
}

func writeWavHeader(file *os.File) {
	header := wavHeader{
		ChunkID:       [4]byte{'R', 'I', 'F', 'F'},
		ChunkSize:     36 + uint32(sampleRate*duration*channels*bitDepth/8),
		Format:        [4]byte{'W', 'A', 'V', 'E'},
		Subchunk1ID:   [4]byte{'f', 'm', 't', ' '},
		Subchunk1Size: 16,
		AudioFormat:   1,
		NumChannels:   uint16(channels),
		SampleRate:    uint32(sampleRate),
		ByteRate:      uint32(sampleRate * channels * bitDepth / 8),
		BlockAlign:    uint16(channels * bitDepth / 8),
		BitsPerSample: uint16(bitDepth),
		Subchunk2ID:   [4]byte{'d', 'a', 't', 'a'},
		Subchunk2Size: uint32(sampleRate * duration * channels * bitDepth / 8),
	}

	binary.Write(file, binary.LittleEndian, header)
}

func writeSample(file *os.File, sample float32) {
	// float32 -> int16変換
	sampleInt16 := int16(sample * (1<<15 - 1))
	binary.Write(file, binary.LittleEndian, sampleInt16)
}
