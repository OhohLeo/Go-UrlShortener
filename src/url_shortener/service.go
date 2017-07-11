// This file is part of Go-UrlShortener.
//
// Go-UrlShortener is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// Go-UrlShortener is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with Go-UrlShortener.  If not, see <http://www.gnu.org/licenses/>.
//
// Authored by OhohLeo
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

const (
	STATUS_UNKNOWN  = -1
	STATUS_STARTING = 0
	STATUS_UP       = 1
	STATUS_STOPPED  = 2
)

type Service struct {
	Dst  string `json:"-"`
	Name string `json:"service"`
	hash string
}

type ServiceConfig struct {
	Port string `json:"port"`
}

type ServiceInstance struct {
	URI    string `json:"uri"`
	Status int    `json:"status"`
	Hash   string `json:"hash"`
}

type ServiceRsp struct {
	Configuration ServiceConfig   `json:"configuration"`
	Instance      ServiceInstance `json:"instance"`
	Registered    bool            `json:"registered"`
}

// Register permet l'enregistrement auprès du PLM
func (s *Service) Register() (ip string, port int, err error) {

	service, err := json.Marshal(s)
	if err != nil {
		return
	}

	// Se déclare auprès du gestionnaire de micro-services
	rsp, err := http.Post(s.Dst+"/register",
		"application/json", bytes.NewBuffer(service))
	if err != nil {
		return
	}

	if rsp.StatusCode != http.StatusOK {
		err = fmt.Errorf("unexpected rsp status %d", rsp.Status)
		return
	}

	defer rsp.Body.Close()

	// Décodage de la réponse
	var serviceRsp ServiceRsp
	err = json.NewDecoder(rsp.Body).Decode(&serviceRsp)
	if err != nil {
		return
	}

	port, err = strconv.Atoi(serviceRsp.Configuration.Port)
	if err != nil {
		return
	}

	// Récupération de l'ip/port spécifié
	if serviceRsp.Instance.URI != "NA" {
		ip = serviceRsp.Instance.URI
	}

	s.hash = serviceRsp.Instance.Hash
	return
}

type Update struct {
	Name   string `json:"service"`
	Status int    `json:"status"`
	Hash   string `json:"hash"`
}

type UpdateRsp struct {
	Uri    string `json:"uri"`
	Status int    `json:"status"`
	Hash   string `json:"hash"`
}

// Update permet l'envoi du status au PLM
func (s *Service) Update(status int) error {

	update := &Update{
		Name:   s.Name,
		Status: status,
		Hash:   s.hash,
	}

	data, err := json.Marshal(update)
	if err != nil {
		return err
	}

	// Création de la requête
	client := &http.Client{}
	req, err := http.NewRequest(
		http.MethodPut, s.Dst+"/status", bytes.NewBuffer(data))
	if err != nil {
		return err
	}

	// Envoie de la requête
	rsp, err := client.Do(req)
	if err != nil {
		return err
	}

	if rsp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected rsp status %d", rsp.Status)
	}

	defer rsp.Body.Close()

	// Décodage de la réponse
	var updateRsp UpdateRsp
	err = json.NewDecoder(rsp.Body).Decode(&updateRsp)
	if err != nil {
		return err
	}

	return nil
}
