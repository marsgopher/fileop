package upyun

import (
	"bufio"
	"os"
	"testing"
	"time"

	"github.com/marsgopher/fileop"
	"github.com/stretchr/testify/suite"
)

type WrapTestSuite struct {
	suite.Suite
	w *Client
}

func (s *WrapTestSuite) SetupSuite() {
	_, ok := os.LookupEnv("UPYUN_BUCKET")
	if ok {
		cfg := Config{
			Bucket:    os.Getenv("UPYUN_BUCKET"),
			Operator:  os.Getenv("UPYUN_OPERATOR"),
			Password:  os.Getenv("UPYUN_PASSWORD"),
			UserAgent: "toolman-migu-test",
		}
		w, err := New(cfg)
		s.Assert().NoError(err)
		s.w = w
	}
}

func TestWrapSuite(t *testing.T) {
	suite.Run(t, &WrapTestSuite{})
}

func (s *WrapTestSuite) TestOpen() {
	if s.w == nil {
		s.T().Skip("require env UPYUN_BUCKET for test")
	}

	path := "/migu/2021-12-04/v6.kwaicdn.com/02/05/202112040205-v6.kwaicdn.com-sep0.gz"

	rd, err := fileop.NewFileReader(s.w, path, fileop.GZIP)
	s.Assert().NoError(err)

	start := time.Now()
	cntLine := 0
	scan := bufio.NewScanner(rd)
	for scan.Scan() {
		//s.T().Log(scan.Text())
		cntLine++
	}
	s.Assert().NoError(scan.Err())

	s.T().Log("line:", cntLine, ", cost:", time.Since(start))
}

func (s *WrapTestSuite) TestReaddirnames() {
	if s.w == nil {
		s.T().Skip("require env UPYUN_BUCKET for test")
	}

	dirPath := "/migu/2021-12-04/v6.kwaicdn.com/02/05/"

	start := time.Now()
	cntLine := 0
	fs, err := s.w.Readdirnames(dirPath, 0)
	s.Assert().NoError(err)

	for _, f := range fs {
		s.T().Log(f)
		cntLine++
	}

	s.T().Log("line:", cntLine, ", cost:", time.Since(start))
}
