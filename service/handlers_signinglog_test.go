// -*- Mode: Go; indent-tabs-mode: t -*-

/*
 * Copyright (C) 2016-2017 Canonical Ltd
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License version 3 as
 * published by the Free Software Foundation.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 *
 */

package service

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSigningLogListHandler(t *testing.T) {
	// Mock the database
	Environ = &Env{DB: &mockDB{}}

	response, _ := sendSigningLogRequest(t, "GET", "/1.0/signinglog", nil)
	if len(response.SigningLog) != 10 {
		t.Errorf("Expected 10 signing logs, got: %d", len(response.SigningLog))
	}
}

func TestSigningLogListHandlerWithParam(t *testing.T) {
	// Mock the database
	Environ = &Env{DB: &mockDB{}}

	response, _ := sendSigningLogRequest(t, "GET", "/1.0/signinglog?fromID=5", nil)
	if len(response.SigningLog) != 4 {
		t.Errorf("Expected 4 signing logs, got: %d", len(response.SigningLog))
	}
}

func TestSigningLogListHandlerBadParam(t *testing.T) {
	// Mock the database
	Environ = &Env{DB: &mockDB{}}

	response, _ := sendSigningLogRequest(t, "GET", "/1.0/signinglog?fromID=bad", nil)
	if len(response.SigningLog) != 10 {
		t.Errorf("Expected 10 signing logs, got: %d", len(response.SigningLog))
	}
}

func TestSigningLogListHandlerError(t *testing.T) {
	// Mock the database
	Environ = &Env{DB: &errorMockDB{}}

	sendSigningLogRequestExpectError(t, "GET", "/1.0/signinglog", nil)
}

func TestSigningLogDeleteHandler(t *testing.T) {
	// Mock the database
	Environ = &Env{DB: &mockDB{}}

	// Delete a signing log
	data := "{}"
	sendSigningLogRequest(t, "DELETE", "/1.0/signinglog/1", bytes.NewBufferString(data))
}

func TestSigningLogDeleteHandlerWrongID(t *testing.T) {
	// Mock the database
	Environ = &Env{DB: &mockDB{}}

	// Delete a signing log
	data := "{}"
	sendSigningLogRequestExpectError(t, "DELETE", "/1.0/signinglog/22", bytes.NewBufferString(data))
}

func TestSigningLogDeleteHandlerError(t *testing.T) {
	// Mock the database
	Environ = &Env{DB: &errorMockDB{}}

	// Delete a signing log
	data := "{}"
	sendSigningLogRequestExpectError(t, "DELETE", "/1.0/signinglog/1", bytes.NewBufferString(data))
}

func TestSigningLogDeleteHandlerBadID(t *testing.T) {
	// Mock the database
	Environ = &Env{DB: &errorMockDB{}}

	// Delete a signing log
	data := "{}"
	sendSigningLogRequestExpectError(t, "DELETE", "/1.0/signinglog/99999999999999999999999999999999999999999999999", bytes.NewBufferString(data))
}

func sendSigningLogRequest(t *testing.T, method, url string, data io.Reader) (SigningLogResponse, error) {
	w := httptest.NewRecorder()
	r, _ := http.NewRequest(method, url, data)
	AdminRouter(Environ).ServeHTTP(w, r)

	// Check the JSON response
	result := SigningLogResponse{}
	err := json.NewDecoder(w.Body).Decode(&result)
	if err != nil {
		t.Errorf("Error decoding the signing log response: %v", err)
	}
	if !result.Success {
		t.Error("Expected success, got error")
	}

	return result, err
}

func sendSigningLogRequestExpectError(t *testing.T, method, url string, data io.Reader) (SigningLogResponse, error) {
	w := httptest.NewRecorder()
	r, _ := http.NewRequest(method, url, data)
	AdminRouter(Environ).ServeHTTP(w, r)

	// Check the JSON response
	result := SigningLogResponse{}
	err := json.NewDecoder(w.Body).Decode(&result)
	if err != nil {
		t.Errorf("Error decoding the signing log response: %v", err)
	}
	if result.Success {
		t.Error("Expected error, got success")
	}

	return result, err
}
