# Coll

[![Build Status](https://www.travis-ci.org/mylxsw/coll.svg?branch=master)](https://www.travis-ci.org/mylxsw/coll)
[![Coverage Status](https://coveralls.io/repos/github/mylxsw/coll/badge.svg?branch=master)](https://coveralls.io/github/mylxsw/coll?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/mylxsw/coll)](https://goreportcard.com/report/github.com/mylxsw/coll)
[![codecov](https://codecov.io/gh/mylxsw/coll/branch/master/graph/badge.svg)](https://codecov.io/gh/mylxsw/coll)
[![Sourcegraph](https://sourcegraph.com/github.com/mylxsw/coll/-/badge.svg)](https://sourcegraph.com/github.com/mylxsw/coll?badge)
[![GitHub](https://img.shields.io/github/license/mylxsw/coll.svg)](https://github.com/mylxsw/coll)



Coll is a collection library for Go.

    cc := coll.MustNew(testMapData)
    collectionWithoutEmpty := cc.Filter(func(value string) bool {
        return value != ""
    }).Filter(func(value string, key string) bool {
        return key != ""
    })
    collectionWithoutEmpty.Each(func(value, key string) {
        if value == "" || key == "" {
            t.Errorf("test failed: %s=>%s", key, value)
        }
    })