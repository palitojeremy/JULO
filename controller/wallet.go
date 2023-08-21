package controller

import (
	"mw-project/middleware"
	"mw-project/model"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

func GetWallet(context *gin.Context) {
	//Token authenticate
	claims := jwt.StandardClaims{}
	headerToken := context.Request.Header["Authorization"][0]
	if headerToken == "" {
		context.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Access denied",
		})
		return
	}
	if err := middleware.AuthenticateToken(middleware.SplitToken(headerToken)); err != nil {
		context.JSON(http.StatusMethodNotAllowed, gin.H{
			"status":  "error",
			"message": "token expired or invalid",
		})
		return
	}
	jwt.ParseWithClaims(middleware.SplitToken(headerToken), &claims, func(token *jwt.Token) (interface{}, error) {
    return []byte(os.Getenv("JWT_SECRET")), nil
	})
	//Token authenticate
	userId := claims.Issuer
	var wallet model.Wallet
	model.DB.Where("owned_by = ?", userId).First(&wallet)

	if wallet.Status == "disabled"{
		context.JSON(http.StatusMethodNotAllowed, gin.H{
			"status":  "fail",
			"message": "Wallet Disabled",
		})
		return
	}
	context.JSON(http.StatusOK, gin.H{"data": wallet, "status":"success"})
}

func EnableWallet(context *gin.Context) {
	//Token authenticate
	claims := jwt.StandardClaims{}
	headerToken := context.Request.Header["Authorization"][0]
	if headerToken == "" {
		context.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Access denied",
		})
		return
	}
	if err := middleware.AuthenticateToken(middleware.SplitToken(headerToken)); err != nil {
		context.JSON(http.StatusMethodNotAllowed, gin.H{
			"status":  "error",
			"message": "token expired or invalid",
		})
		return
	}
	jwt.ParseWithClaims(middleware.SplitToken(headerToken), &claims, func(token *jwt.Token) (interface{}, error) {
    return []byte(os.Getenv("JWT_SECRET")), nil
	})
	//Token authenticate
	userId := claims.Issuer
	var wallet model.Wallet
	model.DB.Where("owned_by = ?", userId).First(&wallet)
	if(wallet.Status == "enabled"){
		context.JSON(http.StatusMethodNotAllowed, gin.H{
			"status":  "fail",
			"message": "Already Enabled",
		})
		return
	}
	model.DB.Model(&wallet).Updates(model.Wallet{Status: "enabled", EnabledAt: time.Now()})
	context.JSON(http.StatusOK, gin.H{"data": wallet, "status": "success"})
	return
}

func CreateWallet(context *gin.Context)  {
	  // Validate input
		var input model.CreateWalletInput
		if err := context.ShouldBindJSON(&input); err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"iss": input.Id,
			"exp": time.Now().Add(time.Hour*24).Unix(),
		})
	
		tokenString, _ := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
		// Create book
		wallet := model.Wallet{Id: uuid.New().String(), OwnedBy: input.Id, Status: "disabled"}
		model.DB.Create(&wallet)
	
		context.JSON(http.StatusOK, gin.H{"data": tokenString, "status": "success"})
		return
}

func GetTransactions(context *gin.Context) {
	//Token authenticate
	claims := jwt.StandardClaims{}
	headerToken := context.Request.Header["Authorization"][0]
	if headerToken == "" {
		context.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Access denied",
		})
		return
	}
	if err := middleware.AuthenticateToken(middleware.SplitToken(headerToken)); err != nil {
		context.JSON(http.StatusMethodNotAllowed, gin.H{
			"status":  "error",
			"message": "token expired or invalid",
		})
		return
	}
	jwt.ParseWithClaims(middleware.SplitToken(headerToken), &claims, func(token *jwt.Token) (interface{}, error) {
    return []byte(os.Getenv("JWT_SECRET")), nil
	})
	//Token authenticate
	userId := claims.Issuer
	var wallet model.Wallet
	model.DB.Where("owned_by = ?", userId).First(&wallet)

	if wallet.Status == "disabled"{
		context.JSON(http.StatusMethodNotAllowed, gin.H{
			"status":  "fail",
			"message": "Wallet Disabled",
		})
		return
	}
	
	var transaction []model.Transaction
	model.DB.Where("transacted_by = ?", userId).Find(&transaction)
	context.JSON(http.StatusOK, gin.H{"data": transaction, "status":"success"})
}

func DepostitWallet(context *gin.Context)  {
	var input model.TransactionInput
	if err := context.ShouldBindJSON(&input); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	//Token authenticate
	claims := jwt.StandardClaims{}
	headerToken := context.Request.Header["Authorization"][0]
	if headerToken == "" {
		context.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Access denied",
		})
		return
	}
	if err := middleware.AuthenticateToken(middleware.SplitToken(headerToken)); err != nil {
		context.JSON(http.StatusMethodNotAllowed, gin.H{
			"status":  "error",
			"message": "token expired or invalid",
		})
		return
	}
	jwt.ParseWithClaims(middleware.SplitToken(headerToken), &claims, func(token *jwt.Token) (interface{}, error) {
    return []byte(os.Getenv("JWT_SECRET")), nil
	})
	//Token authenticate
	userId := claims.Issuer
	var wallet model.Wallet
	model.DB.Where("owned_by = ?", userId).First(&wallet)

	deposit := model.Transaction{
		Id						: uuid.New().String(),
		Type					: "deposit",
		Status				: "success",
		TransactedBy  : userId,
		TransactedAt 	: time.Now(),
		Amount				: input.Amount,
		ReferenceId 	: input.ReferenceId,
	}

	model.DB.Create(&deposit)

	if wallet.Status == "disabled"{
		model.DB.Model(&deposit).Update("status", "failed")
		context.JSON(http.StatusMethodNotAllowed, gin.H{
			"status":  "fail",
			"message": "Wallet Disabled",
		})
		return
	}

	model.DB.Model(&wallet).Update("balance", wallet.Balance + deposit.Amount)
	context.JSON(http.StatusOK, gin.H{"data": deposit, "status": "success"})
	return
}

func WithdrawWallet(context *gin.Context)  {
	var input model.TransactionInput
	if err := context.ShouldBindJSON(&input); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	//Token authenticate
	claims := jwt.StandardClaims{}
	headerToken := context.Request.Header["Authorization"][0]
	if headerToken == "" {
		context.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Access denied",
		})
		return
	}
	if err := middleware.AuthenticateToken(middleware.SplitToken(headerToken)); err != nil {
		context.JSON(http.StatusMethodNotAllowed, gin.H{
			"status":  "error",
			"message": "token expired or invalid",
		})
		return
	}
	jwt.ParseWithClaims(middleware.SplitToken(headerToken), &claims, func(token *jwt.Token) (interface{}, error) {
    return []byte(os.Getenv("JWT_SECRET")), nil
	})
	//Token authenticate
	userId := claims.Issuer
	var wallet model.Wallet
	model.DB.Where("owned_by = ?", userId).First(&wallet)

	withdraw := model.Transaction{
		Id						: uuid.New().String(),
		Type					: "withdraw",
		Status				: "success",
		TransactedBy  : userId,
		TransactedAt 	: time.Now(),
		Amount				: input.Amount,
		ReferenceId 	: input.ReferenceId,
	}

	model.DB.Create(&withdraw)

	if (wallet.Status == "disabled"){
		model.DB.Model(&withdraw).Update("status", "failed")
		context.JSON(http.StatusMethodNotAllowed, gin.H{
			"status":  "fail",
			"message": "Wallet Disabled",
		})
		return
	}

	if (wallet.Balance < withdraw.Amount){
		model.DB.Model(&withdraw).Update("status", "failed")
		context.JSON(http.StatusMethodNotAllowed, gin.H{
			"status":  "fail",
			"message": "Balance Not Enough",
		})
		return
	}

	model.DB.Model(&wallet).Update("balance", wallet.Balance - withdraw.Amount)
	context.JSON(http.StatusOK, gin.H{"data": withdraw, "status": "success"})
	return
}

func DisableWallet(context *gin.Context) {
	var input model.DisableWalletInput
	if err := context.ShouldBindJSON(&input); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	//Token authenticate
	claims := jwt.StandardClaims{}
	headerToken := context.Request.Header["Authorization"][0]
	if headerToken == "" {
		context.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Access denied",
		})
		return
	}
	if err := middleware.AuthenticateToken(middleware.SplitToken(headerToken)); err != nil {
		context.JSON(http.StatusMethodNotAllowed, gin.H{
			"status":  "error",
			"message": "token expired or invalid",
		})
		return
	}
	jwt.ParseWithClaims(middleware.SplitToken(headerToken), &claims, func(token *jwt.Token) (interface{}, error) {
    return []byte(os.Getenv("JWT_SECRET")), nil
	})
	//Token authenticate
	userId := claims.Issuer
	var wallet model.Wallet
	model.DB.Where("owned_by = ?", userId).First(&wallet)
	if(input.IsDisabled){
		model.DB.Model(&wallet).Update("status", "disabled")
		context.JSON(http.StatusOK, gin.H{"data": wallet, "status": "success"})
		return
	}
	model.DB.Model(&wallet).Update("status", "disabled")
	context.JSON(http.StatusOK, gin.H{"data": "failed to disable wallet", "status": "error"})
	return
}