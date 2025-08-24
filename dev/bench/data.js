window.BENCHMARK_DATA = {
  "lastUpdate": 1756067871644,
  "repoUrl": "https://github.com/tomhoffer/darwinium",
  "entries": {
    "Go Benchmarks": [
      {
        "commit": {
          "author": {
            "email": "tomas.hoffer@gmail.com",
            "name": "tomhoffer",
            "username": "tomhoffer"
          },
          "committer": {
            "email": "tomas.hoffer@gmail.com",
            "name": "tomhoffer",
            "username": "tomhoffer"
          },
          "distinct": true,
          "id": "1097374b1af0c3116eb9f5de6c55a927756b931e",
          "message": "feat(ci): Add github action for regression perf benchmark",
          "timestamp": "2025-08-24T22:07:59+02:00",
          "tree_id": "d492b31d2bcfb48a401809dec8176f5dc0a92a16",
          "url": "https://github.com/tomhoffer/darwinium/commit/1097374b1af0c3116eb9f5de6c55a927756b931e"
        },
        "date": 1756066857794,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkExecutor_PerformMutation_DefaultWorkers",
            "value": 6256211,
            "unit": "ns/op\t  720272 B/op\t   20004 allocs/op",
            "extra": "194 times\n2 procs"
          },
          {
            "name": "BenchmarkExecutor_PerformMutation_DefaultWorkers - ns/op",
            "value": 6256211,
            "unit": "ns/op",
            "extra": "194 times\n2 procs"
          },
          {
            "name": "BenchmarkExecutor_PerformMutation_DefaultWorkers - B/op",
            "value": 720272,
            "unit": "B/op",
            "extra": "194 times\n2 procs"
          },
          {
            "name": "BenchmarkExecutor_PerformMutation_DefaultWorkers - allocs/op",
            "value": 20004,
            "unit": "allocs/op",
            "extra": "194 times\n2 procs"
          },
          {
            "name": "BenchmarkExecutor_PerformMutation_UnlimitedWorkers",
            "value": 3433153,
            "unit": "ns/op\t  723233 B/op\t   20009 allocs/op",
            "extra": "363 times\n2 procs"
          },
          {
            "name": "BenchmarkExecutor_PerformMutation_UnlimitedWorkers - ns/op",
            "value": 3433153,
            "unit": "ns/op",
            "extra": "363 times\n2 procs"
          },
          {
            "name": "BenchmarkExecutor_PerformMutation_UnlimitedWorkers - B/op",
            "value": 723233,
            "unit": "B/op",
            "extra": "363 times\n2 procs"
          },
          {
            "name": "BenchmarkExecutor_PerformMutation_UnlimitedWorkers - allocs/op",
            "value": 20009,
            "unit": "allocs/op",
            "extra": "363 times\n2 procs"
          },
          {
            "name": "BenchmarkTournamentSelector_Select",
            "value": 186327,
            "unit": "ns/op\t  448801 B/op\t    1002 allocs/op",
            "extra": "6042 times\n2 procs"
          },
          {
            "name": "BenchmarkTournamentSelector_Select - ns/op",
            "value": 186327,
            "unit": "ns/op",
            "extra": "6042 times\n2 procs"
          },
          {
            "name": "BenchmarkTournamentSelector_Select - B/op",
            "value": 448801,
            "unit": "B/op",
            "extra": "6042 times\n2 procs"
          },
          {
            "name": "BenchmarkTournamentSelector_Select - allocs/op",
            "value": 1002,
            "unit": "allocs/op",
            "extra": "6042 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "tomas.hoffer@gmail.com",
            "name": "tomhoffer",
            "username": "tomhoffer"
          },
          "committer": {
            "email": "tomas.hoffer@gmail.com",
            "name": "tomhoffer",
            "username": "tomhoffer"
          },
          "distinct": true,
          "id": "c551613e82804b088177d3a55379a9008a47970f",
          "message": "feat(executor): Do not show progress bar when running tests",
          "timestamp": "2025-08-24T22:37:10+02:00",
          "tree_id": "6264349c80657f61b15afee625d8917bc6f624dd",
          "url": "https://github.com/tomhoffer/darwinium/commit/c551613e82804b088177d3a55379a9008a47970f"
        },
        "date": 1756067871352,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkExecutor_Loop_DefaultWorkers",
            "value": 34844737,
            "unit": "ns/op\t20423924 B/op\t   62177 allocs/op",
            "extra": "34 times\n4 procs"
          },
          {
            "name": "BenchmarkExecutor_Loop_DefaultWorkers - ns/op",
            "value": 34844737,
            "unit": "ns/op",
            "extra": "34 times\n4 procs"
          },
          {
            "name": "BenchmarkExecutor_Loop_DefaultWorkers - B/op",
            "value": 20423924,
            "unit": "B/op",
            "extra": "34 times\n4 procs"
          },
          {
            "name": "BenchmarkExecutor_Loop_DefaultWorkers - allocs/op",
            "value": 62177,
            "unit": "allocs/op",
            "extra": "34 times\n4 procs"
          },
          {
            "name": "BenchmarkExecutor_Loop_UnlimitedWorkers",
            "value": 16400988,
            "unit": "ns/op\t20422684 B/op\t   62157 allocs/op",
            "extra": "61 times\n4 procs"
          },
          {
            "name": "BenchmarkExecutor_Loop_UnlimitedWorkers - ns/op",
            "value": 16400988,
            "unit": "ns/op",
            "extra": "61 times\n4 procs"
          },
          {
            "name": "BenchmarkExecutor_Loop_UnlimitedWorkers - B/op",
            "value": 20422684,
            "unit": "B/op",
            "extra": "61 times\n4 procs"
          },
          {
            "name": "BenchmarkExecutor_Loop_UnlimitedWorkers - allocs/op",
            "value": 62157,
            "unit": "allocs/op",
            "extra": "61 times\n4 procs"
          },
          {
            "name": "BenchmarkExecutor_PerformMutation_DefaultWorkers",
            "value": 6204521,
            "unit": "ns/op\t  720272 B/op\t   20004 allocs/op",
            "extra": "193 times\n4 procs"
          },
          {
            "name": "BenchmarkExecutor_PerformMutation_DefaultWorkers - ns/op",
            "value": 6204521,
            "unit": "ns/op",
            "extra": "193 times\n4 procs"
          },
          {
            "name": "BenchmarkExecutor_PerformMutation_DefaultWorkers - B/op",
            "value": 720272,
            "unit": "B/op",
            "extra": "193 times\n4 procs"
          },
          {
            "name": "BenchmarkExecutor_PerformMutation_DefaultWorkers - allocs/op",
            "value": 20004,
            "unit": "allocs/op",
            "extra": "193 times\n4 procs"
          },
          {
            "name": "BenchmarkExecutor_PerformMutation_UnlimitedWorkers",
            "value": 2677252,
            "unit": "ns/op\t  724965 B/op\t   20013 allocs/op",
            "extra": "441 times\n4 procs"
          },
          {
            "name": "BenchmarkExecutor_PerformMutation_UnlimitedWorkers - ns/op",
            "value": 2677252,
            "unit": "ns/op",
            "extra": "441 times\n4 procs"
          },
          {
            "name": "BenchmarkExecutor_PerformMutation_UnlimitedWorkers - B/op",
            "value": 724965,
            "unit": "B/op",
            "extra": "441 times\n4 procs"
          },
          {
            "name": "BenchmarkExecutor_PerformMutation_UnlimitedWorkers - allocs/op",
            "value": 20013,
            "unit": "allocs/op",
            "extra": "441 times\n4 procs"
          },
          {
            "name": "BenchmarkTournamentSelector_Select",
            "value": 169557,
            "unit": "ns/op\t  448801 B/op\t    1002 allocs/op",
            "extra": "7101 times\n4 procs"
          },
          {
            "name": "BenchmarkTournamentSelector_Select - ns/op",
            "value": 169557,
            "unit": "ns/op",
            "extra": "7101 times\n4 procs"
          },
          {
            "name": "BenchmarkTournamentSelector_Select - B/op",
            "value": 448801,
            "unit": "B/op",
            "extra": "7101 times\n4 procs"
          },
          {
            "name": "BenchmarkTournamentSelector_Select - allocs/op",
            "value": 1002,
            "unit": "allocs/op",
            "extra": "7101 times\n4 procs"
          }
        ]
      }
    ]
  }
}