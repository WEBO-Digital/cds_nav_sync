package hashrecs

import (
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	filesystem "nav_sync/mods/ahelpers/file_system"
	data_parser "nav_sync/mods/ahelpers/parser"
	"nav_sync/utils"

	"golang.org/x/crypto/sha3"
)

func (hashrecs *HashRecs) Load() {
	// Get the current working directory
	currentDir, _ := filesystem.GetCurrentWorkingDirectory()
	filepath := currentDir + hashrecs.FilePath + hashrecs.Name + ".json"
	utils.Console("filepath----------------------> ", filepath)
	jsonBytes, err := ioutil.ReadFile(filepath)
	if err != nil {
		utils.Console("hashrecs----------------------> ", err.Error())
		hashrecs.Recs = map[string]HashRec{}
	} else {
		json.Unmarshal(jsonBytes, &hashrecs.Recs)
	}
	utils.Console("hashrecs----------------------> ", hashrecs.Recs)
}

func (hashrecs *HashRecs) Set(key string, rec HashRec) {
	hashrecs.Recs[key] = rec
}

func (hashrecs *HashRecs) Get(key string) HashRec {
	return hashrecs.Recs[key]
}

func (hashrecs *HashRecs) Save() bool {
	response, _ := data_parser.ParseModelToString(hashrecs.Recs)
	filesystem.CleanAndSave(hashrecs.FilePath, hashrecs.Name+".json", response)
	return true
}

func (hashrecs *HashRecs) GetHash(key string) string {
	return hashrecs.Recs[key].Hash
}

func Hash(content string) string {
	// hasher := md5.New()
	// hasher.Write([]byte(content))
	// hashInBytes := hasher.Sum(nil)
	// return hex.EncodeToString(hashInBytes)

	hasher := sha3.New512()
	hasher.Write([]byte(content))
	hashInBytes := hasher.Sum(nil)
	return hex.EncodeToString(hashInBytes)
}
