import logo from './logo.svg';
import './App.css';
import {useState} from 'react'
/// use toast for error handeling

const STATES = {
  'register': 1,
  'login': 2,
  'dashboard': 3,
}

const NOTE_STATE = {
  'none': 0,
  'show': 1,
  'edit': 2,
  'post': 3,
}

function Register() {
  /// username
  /// password
}

function Login() {
  /// username
  /// password
  /// link to go to register
}

function ListOfNotes({updateNoteId, updateNoteState}) {
  /// list of notes title
  /// delete button ?
  /// edit button ?
  /// show button
}

function ShowNote({noteId, noteState}) {
  /// markdown support
  /// close button
  /// delete button ?
  /// edit button ?
}

function EditOrPostNote({onSend}) {
  /// component for editing or posting note
  /// send button
}

function Dashboard() {
  /// ListOfNotes
  /// ShowNote if a note selected
  /// EditOrPostNote if edit or post note
  const [noteId, updateNoteId] =  useState(null);
  const [noteState, updateNoteState] = useState(NOTE_STATE.none)
  const onSend = () => {};
  return (
    <>
      {
        <ListOfNotes
          updateNoteId={updateNoteId}
          updateNoteState={updateNoteState}
        />
      }
      {noteState == NOTE_STATE.show && <ShowNote noteId={noteId}  noteState={noteState} updateNoteState={updateNoteState}/>}
      {
        (noteState == NOTE_STATE.edit || noteState == NOTE_STATE.post) &&
        <EditOrPostNote 
          noteId={noteId}
          noteState={noteState}
          updateNoteState={updateNoteState}
        />
      }
    </>
  )
}


function App() {
  const [state, updateState] = useState(STATES.login);
  return (
    <>
      {state == STATES.login && <Login />}
      {state == STATES.register && <Register />}
      {state == STATES.dashboard && <Dashboard />}
    </>
  );
}

export default App;
