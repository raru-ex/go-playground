package interfacebench_test

import "testing"

type TestInterface interface {
	Sum(a, b int64) int64
	SumNoInline(a, b int64) int64
}

type TestImpl struct{}

func (*TestImpl) Sum(a, b int64) int64 {
	return a + b
}

//go:noinline
func (*TestImpl) SumNoInline(a, b int64) int64 {
	return a + b
}

func NewInterfacePointer() TestInterface {
	return &TestImpl{}
}
func NewEntity() TestImpl {
	return TestImpl{}
}
func NewEntityPointer() *TestImpl {
	return &TestImpl{}
}

func BenchmarkInterface(b *testing.B) {
	ipointer := NewInterfacePointer()
	entity := NewEntity()
	pointer := NewEntityPointer()

	// ダウンキャストが発生してポインタ
	b.Run("InterfacePointer", func(b *testing.B) {
		b.ResetTimer() // タイマーをリセット
		for i := 0; i < b.N; i++ {
			ipointer.Sum(1, 1)
		}
	})

	// キャストを明示
	b.Run("InterfacePointer Cast", func(b *testing.B) {
		b.ResetTimer() // タイマーをリセット
		for i := 0; i < b.N; i++ {
			(ipointer).(*TestImpl).Sum(1, 1)
		}
	})

	// ダウンキャストはなくポインタ
	b.Run("Pointer", func(b *testing.B) {
		b.ResetTimer() // タイマーをリセット
		for i := 0; i < b.N; i++ {
			pointer.Sum(1, 1)
		}
	})

	// ダウンキャストはなく実体
	b.Run("Entity", func(b *testing.B) {
		b.ResetTimer() // タイマーをリセット
		for i := 0; i < b.N; i++ {
			entity.Sum(1, 1)
		}
	})

	// ダウンキャストが発生してポインタ
	b.Run("InterfacePointer no-inline", func(b *testing.B) {
		b.ResetTimer() // タイマーをリセット
		for i := 0; i < b.N; i++ {
			ipointer.SumNoInline(1, 1)
		}
	})

	// キャストを明示
	b.Run("InterfacePointer Cast no-inline", func(b *testing.B) {
		b.ResetTimer() // タイマーをリセット
		for i := 0; i < b.N; i++ {
			(ipointer).(*TestImpl).SumNoInline(1, 1)
		}
	})

	// ダウンキャストはなくポインタ
	b.Run("Pointer no-inline", func(b *testing.B) {
		b.ResetTimer() // タイマーをリセット
		for i := 0; i < b.N; i++ {
			pointer.SumNoInline(1, 1)
		}
	})

	// ダウンキャストはなく実体
	b.Run("Entity no-inline", func(b *testing.B) {
		b.ResetTimer() // タイマーをリセット
		for i := 0; i < b.N; i++ {
			entity.SumNoInline(1, 1)
		}
	})

}
