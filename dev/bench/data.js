window.BENCHMARK_DATA = {
  "lastUpdate": 1699683479150,
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
      },
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
          "id": "5cbc6eacd6d80983472ee33d5a85cd23747d6fe4",
          "message": "docs: Add benchmark instructions to CONTRIBUTING\n\nAdd instructions on running benchmarks and viewing data to\nCONTRIBUTING.md.\n\nSigned-off-by: Ian Lewis <ianmlewis@gmail.com>",
          "timestamp": "2023-11-11T06:11:27Z",
          "tree_id": "1d8a4b0632ccc595482f6ac8560f0f6ce994f232",
          "url": "https://github.com/ianlewis/runeio/commit/5cbc6eacd6d80983472ee33d5a85cd23747d6fe4"
        },
        "date": 1699683478681,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkReadSmall",
            "value": 5113,
            "unit": "ns/op",
            "extra": "2342031 times\n4 procs"
          },
          {
            "name": "BenchmarkReadLarge",
            "value": 327059,
            "unit": "ns/op",
            "extra": "36909 times\n4 procs"
          },
          {
            "name": "BenchmarkPeekSmall",
            "value": 1632,
            "unit": "ns/op",
            "extra": "7353036 times\n4 procs"
          },
          {
            "name": "BenchmarkPeekLarge",
            "value": 104203,
            "unit": "ns/op",
            "extra": "114768 times\n4 procs"
          },
          {
            "name": "BenchmarkDiscardSmall",
            "value": 5055,
            "unit": "ns/op",
            "extra": "2429534 times\n4 procs"
          },
          {
            "name": "BenchmarkDiscardLarge",
            "value": 326277,
            "unit": "ns/op",
            "extra": "37897 times\n4 procs"
          }
        ]
      }
    ]
  }
}