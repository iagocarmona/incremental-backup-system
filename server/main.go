package main

import (
	"fmt"
	"log"
	"net/http"
)

type message struct {
	ID int `json:"id"`
	Type string `json:"type"`
	Message string `json:"message"`
}

