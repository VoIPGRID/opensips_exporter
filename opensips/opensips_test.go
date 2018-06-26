package opensips_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/VoIPGRID/opensips_exporter/internal/mock"
	"github.com/VoIPGRID/opensips_exporter/opensips"
	"golang.org/x/sync/errgroup"
)

func TestGetStatistics(t *testing.T) {
	const fakeStatistic = "core:fake_statistic = 42\n"
	var fakeStatisticObject = opensips.Statistic{
		Name:   "fake_statistic",
		Module: "core",
		Value:  42,
	}
	m, err := mock.New([]byte("200 OK\n"+fakeStatistic), 0)
	if err != nil {
		t.Fatal(err)
	}
	o, err := opensips.New(m.Socket())
	if err != nil {
		t.Fatal(err)
	}
	var g errgroup.Group
	g.Go(func() error {
		statistics, err := o.GetStatistics("fake_statistic")
		if err != nil {
			return err
		}
		if len(statistics) != 1 {
			return fmt.Errorf("expected 1 line from GetStatistics, got %d", len(statistics))
		}
		if statistics["fake_statistic"] != fakeStatisticObject {
			return fmt.Errorf("expected %v, got %v", fakeStatistic, statistics["fake_statistic"])
		}
		return nil
	})
	if err := m.Run(1, time.Now().Add(10*time.Second)); err != nil {
		t.Fatal(err)
	}
	if err := g.Wait(); err != nil {
		t.Fatal(err)
	}
	if err := o.Close(); err != nil {
		t.Fatal(err)
	}
	if err := m.Close(); err != nil {
		t.Fatal(err)
	}
}

func TestConcurrentGetStatistics(t *testing.T) {
	const fakeStatistic = "core:fake_statistic = 42\n"
	var fakeStatisticObject = opensips.Statistic{
		Name:   "fake_statistic",
		Module: "core",
		Value:  42,
	}
	m, err := mock.New([]byte("200 OK\n"+fakeStatistic), 100*time.Millisecond)
	if err != nil {
		t.Fatal(err)
	}
	o, err := opensips.New(m.Socket())
	if err != nil {
		t.Fatal(err)
	}
	var g errgroup.Group
	for i := 0; i < 10; i++ {
		g.Go(func() error {
			statistics, err := o.GetStatistics("fake_statistic")
			if err != nil {
				return err
			}
			if len(statistics) != 1 {
				return fmt.Errorf("expected 1 line from GetStatistics, got %d", len(statistics))
			}
			if statistics["fake_statistic"] != fakeStatisticObject {
				return fmt.Errorf("expected %v, got %v", fakeStatisticObject, statistics["fake_statistic"])
			}
			return nil
		})
	}
	if err := m.Run(10, time.Now().Add(10*time.Second)); err != nil {
		t.Fatal(err)
	}
	if err := g.Wait(); err != nil {
		t.Fatal(err)
	}
	if err := o.Close(); err != nil {
		t.Fatal(err)
	}
	if err := m.Close(); err != nil {
		t.Fatal(err)
	}
}
