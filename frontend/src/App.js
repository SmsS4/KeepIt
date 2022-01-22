import logo from './logo.svg';
import './App.css';
import { Form, Input, Button, Checkbox, Card, Spin, message } from 'antd';
import {useState} from 'react'

import { ToastContainer, toast } from 'react-toastify';
import 'react-toastify/dist/ReactToastify.css';

import ReactMarkdown from 'react-markdown'

import { Login} from './LogReg';

const { TextArea } = Input;

const notify = (msg) => toast(msg);

const BASE = "http://185.18.212.202:8080"
const API = `${BASE}/private/`

const STATES = {
  'login': 2,
  'dashboard': 3,
};

const NOTE_STATE = {
  'none': 0,
  'show': 1,
};

const default_text = "# Hello, *world*!\n## salam2\n### salam3\nhello ***majid***";


function ListOfNotes({updateNoteId, updateNoteState, token}) {
  const [notesList, updateNotesList] = useState([]);

  const newNote = () => {
    netNote(token, default_text, data => {
      console.log(data);
      updateNotesList(x => [...x, {id: data.note_id}]);
    });
  };
  const selectNote = (id) => {
    // TODO: send GET request to server
    const success = true;
    if (success) {
      console.log("select", id);
      updateNoteState(x => NOTE_STATE.show); 
    }
  };
  const deleteNote = (id) => {
    // TODO: send DELETE request to server
    console.log("delete", id);
  };
  
  return (
    <>
      <Button type="primary" onClick={newNote}>ساخت یادداشت جدید</Button>
      <div dir="rtl">فهرست یادداشت‌ها:</div>
      <div>
      {notesList.map(function(d, idx){
        return (
          <>
            <div dir="rtl" key={idx}>
              {d.id}
              <Button type="primary" onClick={() => selectNote(d.id)}>انتخاب</Button>
              <Button type="primary" onClick={() => deleteNote(d.id)}>پاک کردن</Button>
            </div>
          </>
        )
      })}
      </div>
    </>
  )
}

function ShowNote({noteId, noteState, updateNoteState, token}) {
  const [editorText, updateEditorText] = useState();
  
  const onChange = e => {
    updateEditorText(x => e.target.value);
  };

  const saveNote = () => {
    // TODO: send PUT request to server by noteID
  };

  const close = () => updateNoteState(x => NOTE_STATE.none);

  return (
    <>
      <TextArea showCount maxLength={1000} style={{ width: "90%", height: 200 }} onChange={onChange} defaultValue={editorText} />
      <div style={{ width: "90%", height: 200 }}>
        <ReactMarkdown>{editorText}</ReactMarkdown>
      </div>
      <div>
        <Button type="primary" onClick={saveNote}>ذخیره یادداشت</Button>
        <Button type="primary" onClick={close}>بستن یادداشت</Button>
      </div>
    </>
  )
}


function Dashboard({updateState, token, updateToken}) {
  const [noteId, updateNoteId] = useState(null);
  const [noteState, updateNoteState] = useState(NOTE_STATE.none);
  const log_out = () => {
    updateToken(x => null);
    updateState(x => STATES.login);
  };
  return (
    <>
      {
        <ListOfNotes
          updateNoteId={updateNoteId}
          updateNoteState={updateNoteState}
          token={token}
        />
      }
      {noteState === NOTE_STATE.show && <ShowNote token={token} noteId={noteId} noteState={noteState} updateNoteState={updateNoteState}/>}
      <div>
        <Button type="primary" onClick={log_out}>خروج از اکانت</Button>
      </div>
    </>
  )
}

function sendRequest(url, method, token, body, callback, q = '') {
  fetch(API+url+q,
  {
    method: method,
    headers: {
      'content-type': 'application/json',
      'Authorization': 'Bearer ' + token,
    },
    body: JSON.stringify(body),
  }
  )
  .then(response => {
    console.log(response)
    return response.json();
  })
  .then(function (data) {
    console.log(data)
    toast.dismiss()
    if (data["error"]){
      toast.error(data["error"])
    }else {
      toast.success(data["message"])
      callback(data)
    }
  })
  .catch(function (error) {
    console.log(error);
  });
}

function getNotesByUsername(token, username, callback) {
  sendRequest(
    "list_note",
    "POST",
    token,
    {'username': username},
    callback
  )
}

function getMyNotes(token, callback) {
  getNotesByUsername(token, "", callback)
}

function netNote(token, note, callback) {
  sendRequest(
    "new_note",
    "POST",
    token,
    {'note': note},
    callback
  )
}

function updateNote(token, noteId, note, callback) {
  sendRequest(
    "update_note",
    "PUT",
    token,
    {'note_id': noteId, 'note': note},
    callback
  )
}

function getGote(token, noteId, callback) {
  sendRequest(
    "get_note",
    "GET",
    token,
    {},
    callback,
    `?note_id=${noteId}`
  )
}

function deleteNote(token, noteId, callback) {
  sendRequest(
    "delete_note",
    "DELETE",
    token,
    {},
    callback,
    `?note_id=${noteId}`
  )
}



function App() {
  const [state, updateState] = useState(STATES.login);
  const [token, updateToken] = useState(null);
  // updateToken("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NDI4NDg1ODMsInVzZXJuYW1lIjoidGVzdCJ9.Yh_GWeD9bswTsyeNpURMQ07T7bgBSwkQyRY51MQI2hg")
  return (
    <>
      {state === STATES.login && <Login updateState={updateState} updateToken={updateToken}/>}
      {state === STATES.dashboard && <Dashboard updateState={updateState} token={token} updateToken={updateToken}/>}
      <ToastContainer />
    </>
  );
}

export default App;
