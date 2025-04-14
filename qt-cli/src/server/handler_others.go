// Copyright (C) 2025 The Qt Company Ltd.
// SPDX-License-Identifier: LicenseRef-Qt-Commercial OR LGPL-3.0-only

package server

import (
	"net/http"
	"os"
	"qtcli/util"
	"time"

	"github.com/gin-gonic/gin"
)

func shutdownServer(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Server shutting down..."})

	go func() {
		time.Sleep(1 * time.Second)
		util.SendSigTermOrKill(os.Getpid())
	}()
}
