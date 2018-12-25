package main

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/render"
	"math/rand"
	"strconv"
	"sync"
	"time"
)

type Match struct {
	//Id of the host player
	host int64

	//Id of the player who joined the session
	player int64

	//Id of the winner player
	winner int64

	//Flag about whether or not a player is joined the session
	playerFoundFlag bool

	//Match Id
	Id int64

	//Activation time
	activationTime int64

	//Flag about whether or not the match has conculuded
	matchConcludedFlag bool

	//Finishing time
	finishTime int64

	//Cause of the win. False for normal win, true for opponents early request
	causeOfTheWin bool

	//Host Colors
	hostH float64
	hostS float64
	hostB float64

	//Player Colors
	playerH float64
	playerS float64
	playerB float64

}

type MatchMaking struct {
	//Index of the first available match
	index int

	//Length of the Match Slice
	length int

	//Match slice
	randomMatchSlice []Match
}

/**
Constructor for Match object
 */
func NewMatch(hostId int64, H float64, S float64, B float64) Match{

	newMatch := Match{}
	newMatch.host = hostId
	newMatch.player = 0
	newMatch.winner = 0
	newMatch.playerFoundFlag = false
	newMatch.Id = int64(time.Now().UnixNano())
	newMatch.matchConcludedFlag = false
	newMatch.hostH = H
	newMatch.hostS = S
	newMatch.hostB = B

	return newMatch
}

/**
Constructor for MatchMaking object
 */
func NewMatchmaking() MatchMaking{
	newMatchmaking := MatchMaking{}
	newMatchmaking.index = len(newMatchmaking.randomMatchSlice)
	newMatchmaking.length = len(newMatchmaking.randomMatchSlice)

	return newMatchmaking
}

func main() {

	var (
		matchMakingSlice []Match
		matchMakingSliceMutex sync.Mutex

		randomMatches = NewMatchmaking()
		randomMatchIndexMutex sync.Mutex
		randomMatchSliceMutex sync.Mutex

		matchMutex sync.Mutex
		matchMakingMutex sync.Mutex

	)

	r := gin.Default()


	/**
	userId: Hosts Id
	H: Hue value of the user
	S: Saturation value of the user
	B: Brightness value of the user
	Holds the response until a match can be made, then returns enemy player's id and the match id to the host
	 */
	r.GET("hostAMatch/:userId/:H/:S/:B", func(c *gin.Context) {

		host, err := strconv.ParseInt(c.Params.ByName("userId"), 10, 64)

		if err != nil {
			panic(err)
		}

		H, err1 := strconv.ParseFloat(c.Params.ByName("H"), 64)

		if err1 != nil {
			panic(err1)
		}

		S, err2 := strconv.ParseFloat(c.Params.ByName("S"), 64)

		if err2 != nil {
			panic(err2)
		}

		B, err3 := strconv.ParseFloat(c.Params.ByName("B"), 64)

		if err3 != nil {
			panic(err3)
		}

		matchAlreadyExistFlag := false
		var existingMatchIndex = -1

		for currentIndex, currentMatch := range matchMakingSlice{

			if(currentMatch.host == host && !currentMatch.playerFoundFlag){
				existingMatchIndex = currentIndex
				matchAlreadyExistFlag = true
			}
		}

		var matchIndex = 0

		if(matchAlreadyExistFlag){

			matchIndex = existingMatchIndex

		}else {

			matchMakingSliceMutex.Lock()

			matchIndex = len(matchMakingSlice)
			matchMakingSlice = append(matchMakingSlice, NewMatch(host, H, S, B))

			matchMakingSliceMutex.Unlock()
		}


		for(!matchMakingSlice[matchIndex].playerFoundFlag){
			time.Sleep(1 * time.Second)
		}

		c.Render(200, render.JSON{
			Data: map[string]interface{}{
				"MatchId": matchMakingSlice[matchIndex].Id,
				"OpponentId": matchMakingSlice[matchIndex].player,
				"MatchType": 1,
				"ActivationTime" : matchMakingSlice[matchIndex].activationTime,
				"H": matchMakingSlice[matchIndex].playerH,
				"S": matchMakingSlice[matchIndex].playerS,
				"B": matchMakingSlice[matchIndex].playerB,
			},
		})

	})

	/**
	userId: Id of the player who is looking for a match
	H: Hue value of the user
	S: Saturation value of the user
	B: Brightness value of the user
	Holds the response until a match can be made, then returns enemy player's id and the match id to the player
	 */
	r.GET("findARandomMatch/:userId/:H/:S/:B", func(c *gin.Context) {

		user, err := strconv.ParseInt(c.Params.ByName("userId"), 10, 64)

		if err != nil {
			panic(err)
		}

		H, err1 := strconv.ParseFloat(c.Params.ByName("H"), 64)

		if err1 != nil {
			panic(err1)
		}

		S, err2 := strconv.ParseFloat(c.Params.ByName("S"), 64)

		if err2 != nil {
			panic(err2)
		}

		B, err3 := strconv.ParseFloat(c.Params.ByName("B"), 64)

		if err3 != nil {
			panic(err3)
		}

		var ret Match

		randomMatchIndexMutex.Lock()

		if(randomMatches.index == (randomMatches.length - 1)){

			randomMatches.randomMatchSlice[randomMatches.index].player = user
			randomMatches.randomMatchSlice[randomMatches.index].playerFoundFlag = true
			randomMatches.randomMatchSlice[randomMatches.index].playerH = H
			randomMatches.randomMatchSlice[randomMatches.index].playerS = S
			randomMatches.randomMatchSlice[randomMatches.index].playerB = B
			randomMatches.randomMatchSlice[randomMatches.index].activationTime = int64(time.Now().Unix()) + int64(rand.Intn(20)) + int64(20)

			ret = randomMatches.randomMatchSlice[randomMatches.index]
			randomMatches.index++

			randomMatchIndexMutex.Unlock()

			c.Render(200, render.JSON{
				Data: map[string]interface{}{
					"MatchId": ret.Id,
					"OpponentId": ret.host,
					"MatchType": 0,
					"ActivationTime": ret.activationTime,
					"H": ret.hostH,
					"S": ret.hostS,
					"B": ret.hostB,
				},
			})

		}else {

			var matchIndex int = randomMatches.index

			randomMatchSliceMutex.Lock()

			randomMatches.randomMatchSlice = append(randomMatches.randomMatchSlice, NewMatch(user, H, S, B))
			randomMatches.length++

			randomMatchSliceMutex.Unlock()
			randomMatchIndexMutex.Unlock()

			for(!randomMatches.randomMatchSlice[matchIndex].playerFoundFlag){
				time.Sleep(2 * time.Second)
			}

			ret = randomMatches.randomMatchSlice[matchIndex]

			c.Render(200, render.JSON{
				Data: map[string]interface{}{
					"MatchId": ret.Id,
					"OpponentId": ret.player,
					"MatchType": 0,
					"ActivationTime": ret.activationTime,
					"H": ret.playerH,
					"S": ret.playerS,
					"B": ret.playerB,
				},
			})
		}
	})

	/**
	hostId: Id of the host player
	userId: Id of the player who is looking for a match
	H: Hue value of the user
	S: Saturation value of the user
	B: Brightness value of the user
	Holds the response until a match can be made, then returns enemy player's id and the match id to the player
	 */
	r.GET("findAMatch/:hostId/:userId/:H/:S/:B", func(c *gin.Context) {

		host, err := strconv.ParseInt(c.Params.ByName("hostId"), 10, 64)

		if err != nil {
			panic(err)
		}

		user, err := strconv.ParseInt(c.Params.ByName("userId"), 10, 64)

		if err != nil {
			panic(err)
		}

		H, err1 := strconv.ParseFloat(c.Params.ByName("H"), 64)

		if err1 != nil {
			panic(err1)
		}

		S, err2 := strconv.ParseFloat(c.Params.ByName("S"), 64)

		if err2 != nil {
			panic(err2)
		}

		B, err3 := strconv.ParseFloat(c.Params.ByName("B"), 64)

		if err3 != nil {
			panic(err3)
		}


		matchFound := false
		var target int = -1

		for currentIndex,currentMatch := range matchMakingSlice{

			if(currentMatch.host == host){
				target = currentIndex
				matchFound = true
			}
		}

		if(matchFound){

			matchMakingSlice[target].player = user
			matchMakingSlice[target].playerFoundFlag = true
			matchMakingSlice[target].playerH = H
			matchMakingSlice[target].playerS = S
			matchMakingSlice[target].playerB = B
			matchMakingSlice[target].activationTime = int64(time.Now().Unix()) + int64(rand.Intn(20)) + int64(20)

			c.Render(200, render.JSON{
				Data: map[string]interface{}{
					"MatchId": matchMakingSlice[target].Id,
					"OpponentId": matchMakingSlice[target].host,
					"MatchType": 1,
					"ActivationTime": matchMakingSlice[target].activationTime,
					"H": matchMakingSlice[target].hostH,
					"S": matchMakingSlice[target].hostS,
					"B": matchMakingSlice[target].hostB,
				},
			})

		}else {

			c.Render(200, render.JSON{
				Data: map[string]interface{}{
					"MatchId": 0,
					"OpponentId": -1,
					"MatchType": 1,
					"ActivationTime": 0,
					"H": 0,
					"S": 0,
					"B": 0,
				},
			})
		}


	})

	/**
	type: match type 0 for public, 1 for private
	matchId: Id of the target match
	userId: One of the two players in the match
	First access to the specified userId part will win the match and the other will loose, return response accepted message
	If a user tries to access before the activation time, that user will loose the match
	 */
	r.GET("match/:Type/:MatchId/:userId", func(c *gin.Context) {

		mode, err := strconv.ParseInt(c.Params.ByName("Type"), 10, 64)

		if err != nil {
			panic(err)
		}

		id, err := strconv.ParseInt(c.Params.ByName("MatchId"), 10, 64)

		if err != nil {
			panic(err)
		}

		user, err := strconv.ParseInt(c.Params.ByName("userId"), 10, 64)

		if err != nil {
			panic(err)
		}

		var targetIndex int = -1

		switch mode {

		case 0:

			for currentIndex, currentMatch := range randomMatches.randomMatchSlice{

				if(currentMatch.Id == id){
					targetIndex = currentIndex
				}
			}

			if(targetIndex >= 0){

				activation := randomMatches.randomMatchSlice[targetIndex].activationTime
				currentTime := int64(time.Now().Unix())

				host := randomMatches.randomMatchSlice[targetIndex].host
				player := randomMatches.randomMatchSlice[targetIndex].player

				matchMutex.Lock()

				if(randomMatches.randomMatchSlice[targetIndex].winner == 0){

					if(currentTime < activation){

						if(host == user){
							randomMatches.randomMatchSlice[targetIndex].winner = player
						}else {
							randomMatches.randomMatchSlice[targetIndex].winner = host
						}

						randomMatches.randomMatchSlice[targetIndex].causeOfTheWin = true

					}else {

						randomMatches.randomMatchSlice[targetIndex].winner = user
						randomMatches.randomMatchSlice[targetIndex].causeOfTheWin = false

					}

					randomMatches.randomMatchSlice[targetIndex].matchConcludedFlag = true
					randomMatches.randomMatchSlice[targetIndex].finishTime = currentTime
				}

				matchMutex.Unlock()

				c.Render(200, render.JSON{
					Data: map[string]interface{}{
						"Response": "OK",
					},
				})

			}else {

				c.Render(200, render.JSON{
					Data: map[string]interface{}{
						"Response": "ERROR",
					},
				})
			}
		case 1:

			for currentIndex, currentMatch := range matchMakingSlice{

				if(currentMatch.Id == id){
					targetIndex = currentIndex
				}
			}

			if(targetIndex >= 0){

				activation := matchMakingSlice[targetIndex].activationTime
				currentTime := int64(time.Now().Unix())

				host := matchMakingSlice[targetIndex].host
				player := matchMakingSlice[targetIndex].player

				matchMakingMutex.Lock()

				if(matchMakingSlice[targetIndex].winner == 0){

					if(currentTime < activation){

						if(host == user){
							matchMakingSlice[targetIndex].winner = player
						}else {
							matchMakingSlice[targetIndex].winner = host
						}

						matchMakingSlice[targetIndex].causeOfTheWin = true

					}else {

						matchMakingSlice[targetIndex].winner = user
						matchMakingSlice[targetIndex].causeOfTheWin = false

					}

					matchMakingSlice[targetIndex].matchConcludedFlag = true
					matchMakingSlice[targetIndex].finishTime = currentTime
				}

				matchMakingMutex.Unlock()

				c.Render(200, render.JSON{
					Data: map[string]interface{}{
						"Response": "OK",
					},
				})

			}else {

				c.Render(200, render.JSON{
					Data: map[string]interface{}{
						"Response": "ERROR",
					},
				})
			}

		default:

			c.Render(200, render.JSON{
				Data: map[string]interface{}{
					"Response": "ERROR",
				},
			})
		}

	})

	/**
	type: match type 0 for public, 1 for private
	matchId: Id of the target match
	Returns winner of the match as well as some additional information
	 */
	r.GET("match/:Type/:MatchId", func(c *gin.Context) {

		mode, err := strconv.ParseInt(c.Params.ByName("Type"), 10, 64)

		if err != nil {
			panic(err)
		}

		id, err := strconv.ParseInt(c.Params.ByName("MatchId"), 10, 64)

		if err != nil {
			panic(err)
		}

		var targetIndex int = -1

		switch mode {

		case 0:

			for currentIndex, currentMatch := range randomMatches.randomMatchSlice{

				if(currentMatch.Id == id){
					targetIndex = currentIndex
				}
			}

			if(targetIndex >= 0){

				for(!randomMatches.randomMatchSlice[targetIndex].matchConcludedFlag){
					time.Sleep(1 * time.Second)
				}

				c.Render(200, render.JSON{
					Data: map[string]interface{}{
						"Winner": randomMatches.randomMatchSlice[targetIndex].winner,
						"CauseOfTheWin": randomMatches.randomMatchSlice[targetIndex].causeOfTheWin,
						"WinningTime": randomMatches.randomMatchSlice[targetIndex].finishTime,
						"ActivationTime": randomMatches.randomMatchSlice[targetIndex].activationTime,
					},
				})

			}else {

				c.Render(200, render.JSON{
					Data: map[string]interface{}{
						"Winner": -1,
						"CauseOfTheWin": false,
						"WinningTime": 0,
						"ActivationTime": randomMatches.randomMatchSlice[targetIndex].activationTime,
					},
				})
			}
		case 1:

			for currentIndex, currentMatch := range matchMakingSlice{

				if(currentMatch.Id == id){
					targetIndex = currentIndex
				}
			}

			if(targetIndex >= 0){

				for(!matchMakingSlice[targetIndex].matchConcludedFlag){
					time.Sleep(1 * time.Second)
				}

				c.Render(200, render.JSON{
					Data: map[string]interface{}{
						"Winner": matchMakingSlice[targetIndex].winner,
						"CauseOfTheWin": matchMakingSlice[targetIndex].causeOfTheWin,
						"WinningTime": matchMakingSlice[targetIndex].finishTime,
						"ActivationTime": matchMakingSlice[targetIndex].activationTime,
					},
				})

			}else {

				c.Render(200, render.JSON{
					Data: map[string]interface{}{
						"Winner": -1,
						"CauseOfTheWin": false,
						"WinningTime": 0,
						"ActivationTime": matchMakingSlice[targetIndex].activationTime,
					},
				})
			}

		default:

			c.Render(200, render.JSON{
				Data: map[string]interface{}{
					"Winner": -1,
					"CauseOfTheWin": false,
					"WinningTime": 0,
					"ActivationTime": 1,
				},
			})
		}

	})

	r.Run()
}
