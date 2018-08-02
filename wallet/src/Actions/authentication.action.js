
import { B_URL } from './../Constants/constants';
import * as types from './../Constants/action.names';
import { createFile, readFile } from './../Utils/Ethereum';
import { lang } from './../Constants/language';
import axios from 'axios';
const keythereum = require('keythereum');
const fs = window.require('fs');
var async = window.require('async');
const electron = window.require('electron');
const remote = electron.remote;
const SENT_DIR = getUserHome() + '/.sentinel';
var ACCOUNT_ADDR = '';
export const KEYSTORE_FILE = SENT_DIR + '/keystore';

function getUserHome() {
    return remote.process.env[(remote.process.platform === 'win32') ? 'USERPROFILE' : 'HOME'];
}

export function sendError(err) {
    if (err) {
        let error;
        if (typeof err === 'object')
            error = JSON.stringify(err);
        else error = err;
        let data = {
            'os': remote.process.platform + remote.process.arch,
            'account_addr': ACCOUNT_ADDR,
            'error_str': error
        }
        fetch(B_URL + '/logs/error', {
            method: 'POST',
            headers: {
                'Accept': 'application/json',
                'Content-type': 'application/json',
                'Access-Control-Allow-Origin': '*'
            },
            body: JSON.stringify(data)
        }).then(function (res) {
            res.json().then(function (resp) {
            })
        })
    }
}

export function setLanguage(lang) {
    return {
        type: types.LANGUAGE,
        payload: lang
    }
}

export function setComponent(component) {
    return {
        type: types.COMPONENT,
        payload: component
    }
}

export const isOnline = function () {
    try {
        if (window.navigator.onLine) {
            return true
        }
        else {
            return false
        }
    } catch (Err) {
        sendError(Err);
    }
}

export const uploadKeystore = (keystore, cb) => {
    try {
        cb(null, createFile(KEYSTORE_FILE, keystore))
    } catch (Err) {
        sendError(Err);
    }
}

export const createAccount = (password) => {
    try {
        let request = axios({
            url: '/client/account',
            method: 'POST',
            headers: {
                'Access-Control-Allow-Origin': '*',
            },
            data: {
                password: password
            }
        })

        return {
            type: types.CREATE_ACCOUNT,
            payload: request
        };
    } catch (Err) {
        sendError(Err);
    }
}

export function getPrivateKey(password, language, cb) {
    readFile(KEYSTORE_FILE, function (err, data) {
        if (err) cb(err, null);
        else {
            var keystore = JSON.parse(data)
            try {
                var privateKey = keythereum.recover(password, keystore);
                cb(null, privateKey);
            }
            catch (err) {
                cb({ message: lang[language].KeyPassMatch }, null);
            }
        }
    })
}

