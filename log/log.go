// Copyright 2017 The Ethermis Authors
// This file is part of Ethermis.
//
// Ethermis is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// Ethermis is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with Ethermis. If not, see <http://www.gnu.org/licenses/>.

package log

import "github.com/ethereum/go-ethereum/logger/glog"

func init() {
	// mainLogger = logger.("ethermint")
	// sys := logger.NewStdLogSystem(os.Stdout, log.Ldate|log.Ltime, logger.DebugDetailLevel)
	// logger.AddLogSystem(sys)
	glog.SetToStderr(true)
}

func Infof(format string, args ...interface{}) {
	glog.Infof(format, args...)
}

func Warningf(format string, args ...interface{}) {
	glog.Warningf(format, args...)
}

func Errorf(format string, args ...interface{}) {
	glog.Errorf(format, args...)
}

func Fatalf(format string, args ...interface{}) {
	glog.Fatalf(format, args...)
}

func Info(args ...interface{}) {
	glog.Infoln(args...)
}

func Warning(args ...interface{}) {
	glog.Warning(args...)
}

func Error(args ...interface{}) {
	glog.Errorln(args...)
}

func Fatal(args ...interface{}) {
	glog.Fatalln(args...)
}

func Flush() {
	glog.Flush()
}
