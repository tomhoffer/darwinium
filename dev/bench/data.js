window.BENCHMARK_DATA = {
  "lastUpdate": 1756066858456,
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
      }
    ]
  }
}