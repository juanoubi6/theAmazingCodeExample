package address

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"theAmazingCodeExample/app/models"
	"theAmazingCodeExample/app/helpers/googlePlaces"
	"theAmazingCodeExample/app/common"
)

func GetAddresses(c *gin.Context){

	userID := c.MustGet("id").(uint)

	userAddresses,err := models.GetAddressesForUser(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{ "description": "Something went wrong when looking at the addresses", "detail": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{ "description": userAddresses})

}

func AddAddress(c *gin.Context) {

	userID := c.MustGet("id").(uint)

	addressVal := c.PostForm("address")
	floor := c.PostForm("floor")
	apartment := c.PostForm("apartment")
	addressPostalCode := c.PostForm("postal_code")

	//Check for obligatory values
	if addressVal == "" || addressPostalCode == "" {
		c.JSON(http.StatusBadRequest, gin.H{ "description": "Some parameters are missing"})
		return
	}

	//Validate address with google maps api and get the first one that best matches the requested one
	addressName,lat,long,err := googlePlaces.GetValidAddress(addressVal)
	if err != nil {
		if err.Error() == "maps: ZERO_RESULTS - "{
			c.JSON(http.StatusBadRequest, gin.H{ "description": "No result were found for the submitted address", "detail": err.Error()})
			return
		}else {
			c.JSON(http.StatusInternalServerError, gin.H{ "description": err.Error(), "detail": err.Error()})
			return
		}
	}

	//Check postal code is valid. If so, get it's ID
	postalCodeData, found, err := models.GetPostalCodeByCode(addressPostalCode)
	if found == false {
		c.JSON(http.StatusBadRequest, gin.H{ "description": "The address isn't in our address ranges"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{ "description": "Something went wrong", "detail": err.Error()})
		return
	}

	//Get user data
	userData, found, err := models.GetUserById(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{ "description": "Something went wrong", "detail": err.Error()})
		return
	}
	if found == false {
		c.JSON(http.StatusBadRequest, gin.H{ "description": "User not found"})
		return
	}

	//Check if user has any main address
	_,foundMainAddress, err := userData.GetMainAddress()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{ "description": "Something went wrong", "detail": err.Error()})
		return
	}

	//Create new address
	newAddress := models.Address{
		Address:      		addressName,
		Floor:        		floor,
		Apartment:    		apartment,
		MainAddress:  		!foundMainAddress,
		PostalCodeID: 		postalCodeData.ID,
		Latitude:	  		lat,
		Longitude:    		long,
		PostalCode:   		postalCodeData,
		UserID:       		userID,
	}

	if err := newAddress.Save(); err != nil{
		c.JSON(http.StatusInternalServerError, gin.H{ "description": "Something went wrong", "detail": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{ "description": newAddress})

}

func ModifyAddress(c *gin.Context) {

	userID := c.MustGet("id").(uint)

	addressID := c.Param("id")
	addressVal := c.PostForm("address")
	floor := c.PostForm("floor")
	apartment := c.PostForm("apartment")
	addressPostalCode := c.PostForm("postal_code")

	addressIdValue, err := common.StringToUint(addressID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{ "description": "Invalid address ID"})
		return
	}

	//Check for obligatory values
	if addressVal == "" || addressPostalCode == "" {
		c.JSON(http.StatusBadRequest, gin.H{ "description": "Some parameters are missing"})
		return
	}

	//Validate address with google maps api and get the first one that best matches the requested one
	addressName,lat,long,err := googlePlaces.GetValidAddress(addressVal)
	if err != nil {
		if err.Error() == "maps: ZERO_RESULTS - "{
			c.JSON(http.StatusBadRequest, gin.H{ "description": "No result were found for the submitted address", "detail": err.Error()})
			return
		}else {
			c.JSON(http.StatusInternalServerError, gin.H{ "description": "Something went wrong", "detail": err.Error()})
			return
		}
	}

	//Check postal code is valid. If so, get it's ID
	postalCodeData, found, err := models.GetPostalCodeByCode(addressPostalCode)
	if found == false {
		c.JSON(http.StatusBadRequest, gin.H{ "description": "The address isn't in our address ranges", "detail": err.Error()})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{ "description": "Something went wrong", "detail": err.Error()})
		return
	}

	//Get address
	addressData, found, err := models.GetAddressById(addressIdValue)
	if found == false {
		c.JSON(http.StatusBadRequest, gin.H{ "description": "Address not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{ "description": "Something went wrong", "detail": err.Error()})
		return
	}

	//Check you can modify this address
	if addressData.UserID != userID {
		c.JSON(http.StatusUnauthorized, gin.H{ "description": "You can't modify this address"})
		return
	}

	addressData.Address = addressName
	addressData.Floor = floor
	addressData.Apartment = apartment
	addressData.PostalCodeID = postalCodeData.ID
	addressData.Latitude = lat
	addressData.Longitude = long

	if err := addressData.Modify(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{ "description": "Something went wrong", "detail": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{ "description": addressData})

}

func DeleteAddress(c *gin.Context) {

	userID := c.MustGet("id").(uint)
	addressID := c.Param("id")

	addressIdValue, err := common.StringToUint(addressID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{ "description": "Invalid address ID"})
		return
	}

	//Get address
	addressData, found, err := models.GetAddressById(addressIdValue)
	if found == false {
		c.JSON(http.StatusBadRequest, gin.H{ "description": "Address not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{ "description": "Something went wrong", "detail": err.Error()})
		return
	}

	//Check you can delete this address
	if addressData.UserID != userID {
		c.JSON(http.StatusUnauthorized, gin.H{ "description": "You can't modify this address"})
		return
	}

	if err := addressData.Delete(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{ "description": "Something went wrong", "detail": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{ "description": "Address removed successfully"})

}

func MarkAsMain(c *gin.Context) {

	userID := c.MustGet("id").(uint)
	addressID := c.Param("id")

	addressIdValue, err := common.StringToUint(addressID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{ "description": "Invalid address ID"})
		return
	}

	//Get new main address
	newMainAddressData, found, err := models.GetAddressById(addressIdValue)
	if found == false {
		c.JSON(http.StatusBadRequest, gin.H{ "description": "Address not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{ "description": "Something went wrong", "detail": err.Error()})
		return
	}

	//Check you can modify both addresses
	if newMainAddressData.UserID != userID {
		c.JSON(http.StatusUnauthorized, gin.H{ "description": "You can't modify this address"})
		return
	}

	//Get user data
	userData, found, err := models.GetUserById(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{ "description": "Something went wrong", "detail": err.Error()})
		return
	}
	if found == false {
		c.JSON(http.StatusBadRequest, gin.H{ "description": "User not found"})
		return
	}
	
	//Get old main address
	oldMainAddressData, found, err := userData.GetMainAddress()
	if found == true {

		//Check you can modify this address
		if oldMainAddressData.UserID != userID {
			c.JSON(http.StatusUnauthorized, gin.H{ "description": "You can't modify this address"})
			return
		}

		oldMainAddressData.MainAddress = false

		if err := oldMainAddressData.Modify(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{ "description": "Something went wrong", "detail": err.Error()})
			return
		}

	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{ "description": "Something went wrong", "detail": err.Error()})
		return
	}
	
	//Modify new main address
	newMainAddressData.MainAddress = true

	if err := newMainAddressData.Modify(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{ "description": "Something went wrong", "detail": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{ "description": "Main address modified successfully"})

}
