
import { sendError } from '../Actions/authentication.action';
const fs = window.require('fs');

export function createFile(KEYSTORE_FILE, keystore) {
    fs.writeFile(KEYSTORE_FILE, keystore, function (err) {
        if (err) {
            return (err, null);
        } else {
            return KEYSTORE_FILE
        }
    });
}

export function readFile(KEYSTORE_FILE, cb) {
    try {
        fs.readFile(KEYSTORE_FILE, 'utf8', function (err, data) {
            if (err) {
                sendError(err);
                cb(err, null);
            }
            else {
                cb(null, data);
            }
        });
    } catch (Err) {
        sendError(Err);
    }
}