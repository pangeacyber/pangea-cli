package intel

import (
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/pangeacyber/pangea-cli/v2/plugins/utils"
	"github.com/spf13/cobra"
)

var PluginIntelFilePatternReputation = utils.NewPluginWithFilePattern([]string{"file-intel", "v2", "/reputation"}, "hashes", filesToHashes, responseProcess, "file-intel")

type hashFunction func([]byte) string

func HashSHA256(i []byte) string {
	b := sha256.Sum256(i)
	return hex.EncodeToString(b[:])
}

func HashSHA1(i []byte) string {
	b := sha1.Sum(i)
	return hex.EncodeToString(b[:])
}

func responseProcess(filesToValues map[string]string, response map[string]any) error {
	if response == nil {
		return errors.New("got nil response")
	}

	data := response["data"]
	hashesMap, ok := data.(map[string]any)
	if !ok {
		return errors.New("unable to process `data` field")
	}

	output := map[string]any{}

	for k, v := range filesToValues {
		var r any
		if r, ok = hashesMap[v]; !ok {
			continue
		}
		data := map[string]any{}
		data["hash"] = v
		data["result"] = r
		output[k] = data
	}

	b, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		return err
	}
	fmt.Fprintln(os.Stderr, string(b))
	return nil
}

func filesToHashes(cmd *cobra.Command, files []string) (map[string]string, error) {
	flag := cmd.Flag("hash_type")
	if flag == nil {
		return nil, errors.New("`hash_type` flag is not present")
	}
	ht := flag.Value.String()
	if ht == "" {
		ht = "sha256"
		_ = flag.Value.Set("sha256")
		flag.Changed = true
	}

	var hf hashFunction

	switch ht {
	case "sha256":
		hf = HashSHA256
	case "sha1":
		hf = HashSHA1
	default:
		return nil, fmt.Errorf("`%s` not supported on file pattern command", ht)
	}

	fileToHash := map[string]string{}

	for _, file := range files {
		b, err := readAll(file)
		if err != nil {
			return nil, err
		}
		h := hf(b)
		fileToHash[file] = h
	}

	return fileToHash, nil
}

func readAll(filename string) ([]byte, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return io.ReadAll(file)
}
