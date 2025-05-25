module github.com/amyy54/pwlist

go 1.23.5

replace github.com/amyy54/pwlist/internal/reader => ./internal/reader

replace github.com/amyy54/pwlist/internal/cartesian => ./internal/cartesian

replace github.com/amyy54/pwlist/internal/formatter => ./internal/formatter

require github.com/QMUL/ntlmgen v0.0.0-20160211164635-c5fd3399f820 // indirect
