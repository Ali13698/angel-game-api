package main

import (
	"crypto/rand"
	"encoding/binary"
)

func rollDie() int {
	var b [8]byte
	_, _ = rand.Read(b[:])
	return int(binary.LittleEndian.Uint64(b[:])%6) + 1
}

func absolutePos(seat, rel int) int {
	if rel < 0 { return -1 }
	if rel >= MainTrack { return -2 }
	return (SeatStart[seat] + rel) % MainTrack
}

func findKills(s *State, seat, absPos int) []int {
	if absPos < 0 || SafeSquares[absPos] { return nil }
	var killed []int
	for other := 0; other < Seats; other++ {
		if other == seat || s.Seats[other] == "" { continue }
		for j := 0; j < TokensPerSeat; j++ {
			idx := other*TokensPerSeat + j
			pos := s.Tokens[idx]
			if pos < 0 || pos >= MainTrack { continue }
			if absolutePos(other, pos) == absPos { killed = append(killed, idx) }
		}
	}
	return killed
}

func allFinished(s *State, seat int) bool {
	base := seat * TokensPerSeat
	for i := 0; i < TokensPerSeat; i++ {
		if s.Tokens[base+i] != Finish { return false }
	}
	return true
}

func nextTurn(s *State) int {
	n := (s.Turn + 1) % Seats
	for i := 0; i < Seats; i++ {
		if s.Seats[n] != "" { return n }
		n = (n + 1) % Seats
	}
	return n
}
