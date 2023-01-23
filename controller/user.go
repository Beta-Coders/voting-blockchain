package controller

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"votingblockchain/ECC"
	"votingblockchain/container"
	"votingblockchain/model"
)

type UserController struct {
	container container.Container
}

func NewUserController(container container.Container) *UserController {
	return &UserController{container: container}
}

// Login - get user by id
func (controller *UserController) Login(c echo.Context) error {
	user := new(model.User)
	if err := c.Bind(user); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if !controller.checkUserExists(user.Username) {
		return echo.NewHTTPError(http.StatusNotFound, model.Response{Message: "User does not exist"})
	}
	rows, err := controller.container.GetDB().Query("SELECT * FROM voting_users WHERE username=$1", user.Username)
	if err != nil {
		fmt.Println(err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())

	}
	password := user.Password
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&user.Id, &user.Username, &user.Password)
		if err != nil {
			fmt.Println(err)
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
	}
	if !checkPasswordHash(password, user.Password) {
		return echo.NewHTTPError(http.StatusBadRequest, model.Response{Message: "wrong password"})
	}
	return c.JSON(http.StatusOK, model.Response{Message: "success"})
}

// Signup - create user
func (controller *UserController) Signup(c echo.Context) (err error) {
	user := new(model.User)
	if err = c.Bind(user); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	user.Password, _ = hashPassword(user.Password)
	if controller.checkUserExists(user.Username) {
		return echo.NewHTTPError(http.StatusBadRequest, model.Response{Message: "User already exists"})
	}
	rows, err := controller.container.GetDB().Query("INSERT INTO voting_users (id,username,password) VALUES ($1,$2,$3)", user.Id, user.Username, user.Password)
	if err != nil {
		fmt.Println(err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	defer rows.Close()
	//fmt.Println(user)
	pv1, pb1, _, _, signature, signhash := ECC.GenKeys(user.Username)
	pv1Encoded, pb1Encoded := encode(pv1, &pb1)
	//pv2Encoded, pb2Encoded := encode(pk2, &pb2)
	result := model.SignUp{
		Message:   "Signup success",
		PvEncoded: pv1Encoded,
		PbEncoded: pb1Encoded,
		Signature: string(signature),
		SignHash:  string(signhash),
	}
	_, err = controller.container.GetDB().Exec("INSERT INTO admins (public_key,username) VALUES ($1,$2)", pb1Encoded, user.Username)
	//fmt.Println(bytes.Compare(signature, []byte(result.Signature)))
	return c.JSON(http.StatusOK, result)
}

func (controller *UserController) checkUserExists(username string) bool {
	rows, err := controller.container.GetDB().Query("SELECT * FROM voting_users WHERE username=$1", username)
	if err != nil {
		fmt.Println(err)
		return false
	}
	defer rows.Close()
	return rows.Next()
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
func encode(privateKey *ecdsa.PrivateKey, publicKey *ecdsa.PublicKey) (string, string) {
	x509Encoded, _ := x509.MarshalECPrivateKey(privateKey)
	pemEncoded := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: x509Encoded})

	x509EncodedPub, _ := x509.MarshalPKIXPublicKey(publicKey)
	pemEncodedPub := pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: x509EncodedPub})

	return string(pemEncoded), string(pemEncodedPub)
}

func decode(pemEncoded string, pemEncodedPub string) (*ecdsa.PrivateKey, *ecdsa.PublicKey) {
	block, _ := pem.Decode([]byte(pemEncoded))
	x509Encoded := block.Bytes
	privateKey, _ := x509.ParseECPrivateKey(x509Encoded)

	blockPub, _ := pem.Decode([]byte(pemEncodedPub))
	x509EncodedPub := blockPub.Bytes
	genericPublicKey, _ := x509.ParsePKIXPublicKey(x509EncodedPub)
	publicKey := genericPublicKey.(*ecdsa.PublicKey)

	return privateKey, publicKey
}
