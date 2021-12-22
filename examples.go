package main

// func aProgressBar() {
// 	const n = 20
// 	builder := aec.EmptyBuilder

// 	up2 := aec.Up(2)
// 	col := aec.Column(n + 2)
// 	bar := aec.Color8BitF(aec.NewRGB8Bit(64, 255, 64))
// 	label := builder.LightRedF().Underline().With(col).Right(1).ANSI

// 	// for up2
// 	fmt.Println()
// 	fmt.Println()

// 	for i := 0; i <= n; i++ {
// 		fmt.Print(up2)
// 		fmt.Println(label.Apply(fmt.Sprint(i, "/", n)))
// 		fmt.Print("[")
// 		fmt.Print(bar.Apply(strings.Repeat("=", i)))
// 		fmt.Println(col.Apply("]"))
// 		time.Sleep(30 * time.Millisecond)
// 	}
// }
