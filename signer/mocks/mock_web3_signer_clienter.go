// Code generated by mockery. DO NOT EDIT.

package mocks

import (
	context "context"

	common "github.com/ethereum/go-ethereum/common"

	mock "github.com/stretchr/testify/mock"
)

// Web3SignerClienter is an autogenerated mock type for the Web3SignerClienter type
type Web3SignerClienter struct {
	mock.Mock
}

type Web3SignerClienter_Expecter struct {
	mock *mock.Mock
}

func (_m *Web3SignerClienter) EXPECT() *Web3SignerClienter_Expecter {
	return &Web3SignerClienter_Expecter{mock: &_m.Mock}
}

// EthAccounts provides a mock function with given fields: ctx
func (_m *Web3SignerClienter) EthAccounts(ctx context.Context) ([]common.Address, error) {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for EthAccounts")
	}

	var r0 []common.Address
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) ([]common.Address, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) []common.Address); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]common.Address)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Web3SignerClienter_EthAccounts_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'EthAccounts'
type Web3SignerClienter_EthAccounts_Call struct {
	*mock.Call
}

// EthAccounts is a helper method to define mock.On call
//   - ctx context.Context
func (_e *Web3SignerClienter_Expecter) EthAccounts(ctx interface{}) *Web3SignerClienter_EthAccounts_Call {
	return &Web3SignerClienter_EthAccounts_Call{Call: _e.mock.On("EthAccounts", ctx)}
}

func (_c *Web3SignerClienter_EthAccounts_Call) Run(run func(ctx context.Context)) *Web3SignerClienter_EthAccounts_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *Web3SignerClienter_EthAccounts_Call) Return(_a0 []common.Address, _a1 error) *Web3SignerClienter_EthAccounts_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *Web3SignerClienter_EthAccounts_Call) RunAndReturn(run func(context.Context) ([]common.Address, error)) *Web3SignerClienter_EthAccounts_Call {
	_c.Call.Return(run)
	return _c
}

// SignHash provides a mock function with given fields: ctx, address, hashToSign
func (_m *Web3SignerClienter) SignHash(ctx context.Context, address common.Address, hashToSign common.Hash) ([]byte, error) {
	ret := _m.Called(ctx, address, hashToSign)

	if len(ret) == 0 {
		panic("no return value specified for SignHash")
	}

	var r0 []byte
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, common.Address, common.Hash) ([]byte, error)); ok {
		return rf(ctx, address, hashToSign)
	}
	if rf, ok := ret.Get(0).(func(context.Context, common.Address, common.Hash) []byte); ok {
		r0 = rf(ctx, address, hashToSign)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, common.Address, common.Hash) error); ok {
		r1 = rf(ctx, address, hashToSign)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Web3SignerClienter_SignHash_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SignHash'
type Web3SignerClienter_SignHash_Call struct {
	*mock.Call
}

// SignHash is a helper method to define mock.On call
//   - ctx context.Context
//   - address common.Address
//   - hashToSign common.Hash
func (_e *Web3SignerClienter_Expecter) SignHash(ctx interface{}, address interface{}, hashToSign interface{}) *Web3SignerClienter_SignHash_Call {
	return &Web3SignerClienter_SignHash_Call{Call: _e.mock.On("SignHash", ctx, address, hashToSign)}
}

func (_c *Web3SignerClienter_SignHash_Call) Run(run func(ctx context.Context, address common.Address, hashToSign common.Hash)) *Web3SignerClienter_SignHash_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(common.Address), args[2].(common.Hash))
	})
	return _c
}

func (_c *Web3SignerClienter_SignHash_Call) Return(_a0 []byte, _a1 error) *Web3SignerClienter_SignHash_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *Web3SignerClienter_SignHash_Call) RunAndReturn(run func(context.Context, common.Address, common.Hash) ([]byte, error)) *Web3SignerClienter_SignHash_Call {
	_c.Call.Return(run)
	return _c
}

// NewWeb3SignerClienter creates a new instance of Web3SignerClienter. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewWeb3SignerClienter(t interface {
	mock.TestingT
	Cleanup(func())
}) *Web3SignerClienter {
	mock := &Web3SignerClienter{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
