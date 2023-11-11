window.BENCHMARK_DATA = {
  "lastUpdate": 1699682462407,
  "repoUrl": "https://github.com/ianlewis/runeio",
  "entries": {
    "Benchmark": [
      {
        "commit": {
          "author": {
            "email": "ianmlewis@gmail.com",
            "name": "Ian Lewis",
            "username": "ianlewis"
          },
          "committer": {
            "email": "ianmlewis@gmail.com",
            "name": "Ian Lewis",
            "username": "ianlewis"
          },
          "distinct": false,
          "id": "1b97e881c1918bd1ff26fbdea57370478dedba0f",
          "message": "ci: Fix benchmark push on main workflow\n\nSigned-off-by: Ian Lewis <ianmlewis@gmail.com>",
          "timestamp": "2023-11-11T05:55:53Z",
          "tree_id": "8965dd43c42bb9058baf35808068ea3d6cf8fda2",
          "url": "https://github.com/ianlewis/runeio/commit/1b97e881c1918bd1ff26fbdea57370478dedba0f"
        },
        "date": 1699682461852,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkReadSmall",
            "value": 5107,
            "unit": "ns/op",
            "extra": "2345428 times\n4 procs"
          },
          {
            "name": "BenchmarkReadLarge",
            "value": 327158,
            "unit": "ns/op",
            "extra": "36847 times\n4 procs"
          },
          {
            "name": "BenchmarkPeekSmall",
            "value": 1638,
            "unit": "ns/op",
            "extra": "7304431 times\n4 procs"
          },
          {
            "name": "BenchmarkPeekLarge",
            "value": 104074,
            "unit": "ns/op",
            "extra": "115509 times\n4 procs"
          },
          {
            "name": "BenchmarkDiscardSmall",
            "value": 5051,
            "unit": "ns/op",
            "extra": "2433334 times\n4 procs"
          },
          {
            "name": "BenchmarkDiscardLarge",
            "value": 324575,
            "unit": "ns/op",
            "extra": "37951 times\n4 procs"
          }
        ]
      }
    ]
  }
}