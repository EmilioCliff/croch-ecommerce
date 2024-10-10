package handlers

import "github.com/gin-gonic/gin"

func (s *HttpServer) createUser(ctx *gin.Context) {}

func (s *HttpServer) loginUser(ctx *gin.Context) {}

func (s *HttpServer) refreshToken(ctx *gin.Context) {}

func (s *HttpServer) resetPassword(ctx *gin.Context) {}

func (s HttpServer) listUsers(ctx *gin.Context) {}

func (s *HttpServer) updateUserSubscription(ctx *gin.Context) {}

func (s *HttpServer) getUser(ctx *gin.Context) {}
