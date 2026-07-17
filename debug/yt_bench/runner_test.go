package yt_bench

import (
	"testing"
)

func TestRunBenchmarks(t *testing.T) {
	playlistIDs := []string{
		// "PLQY2H8rRoyvzDbLUZkbudP-MFQZwNmU4S",
		// "PLOU2XLYxmsIIuiBfYad6rFYQU_jL2ryal",
		// "PLf7L7Kg8_FNxHATtLwDceyh72QQL9pvpQ",
		// "PLdOKnrf8EcP384Ilxra4UlK9BDJGwawg9",
		// "PLjgj6kdf_snaw8QnlhK5f3DzFDFKDU5f4",
		// "PLNcg_FV9n7qZGfFl2ANI_zISzNp257Lwn",
		// "PLsyeobzWxl7omDoEYrrf3oXvXxa6MPgek",
		// "PLMWaZteqtEaI2Xd7-lnv2hsdMVteY7U1v",
		// "PLHq_wPEVVWy0Vq72bS9wRkw_9Vqgqp0y5",
		// "PLu0W_9lII9agwh1XjRt242xIpHhPT2llg",
		// "PLgUwDviBIf0oF6QL8m22w1hIDC1vJ_BHz",
		// "PLakjEKwjPZDYZMtzRfB26Un77KI2xI_oB" // 100
		// "PLiy0XOfUv4hFH2HbflPOARBXA6qN90mHt" //200
		// "PLtJCksbabIPyfutcI3Kx9mg6h87p2DXzE" //500
		// "PL-CA1f0J88gDsSdoXFz3vIefI8SnBan1s" // 584


		// "PLoaTDHRuxwgxk0WTXRJJJlXUYuHBSrbcF" //1000
		// "PL5fRL6A4m-DFddDPJU5Ugr3NRhxrm5GtR" //1000
		// "PL0dTxWJ6ngUKlBw5eDv7qYylA3xP3Asef" //2198
		"PLdSukIYrTISE", //5000
	}
	RunBenchmarks(playlistIDs)
}
