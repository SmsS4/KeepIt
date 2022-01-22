import { ToastContainer, toast } from 'react-toastify';
import {  Input, Button, Checkbox, Card, Spin, message } from 'antd';
import React from 'react';

const BASE = "http://localhost:8080"
const LOGIN_URL = `${BASE}/pub/login/`
const REGISTER_URL = `${BASE}/pub/register/`

const STATES = {
  'login': 2,
  'dashboard': 3,
};


const appStyle = {
  display: 'flex',
  height: '980px',
  backgroundImage: "url(back.jpeg)",
  backgroundSize: '100% 100%',
  backgroundRepeat: 'no-repeat',
};

const formStyle = {
  transform: "scale(2, 2)",
  marginLeft: '40%',
  marginTop: '28%',
  padding: '10px',
  // border: '1px solid #c9c9c9',
  // borderRadius: '5px',
  // background: '#f5f5f5',
  width: '220px',
  display: 'block'
};



const labelStyle = {
  margin: '10px 0 5px 0',
  fontFamily: 'Arial, Helvetica, sans-serif',
  fontSize: '15px',
};

const inputStyle = {
  margin: '5px 0 10px 0',
  padding: '5px', 
  border: '1px solid #bfbfbf',
  borderRadius: '3px',
  boxSizing: 'border-box',
  width: '100%'
};

const submitStyleHover = {
  
}

const submitStyle = {
  "&:hover": {
    transform: "scale(1.1, 1.1)",
  },
  margin: '10px 0 0 0',
  padding: '7px 10px',
  // border: '1px solid #efffff',
  // borderRadius: '3px',
  // background: '#3085d6',
  width: '100%', 
  fontSize: '15px',
  // color: 'white',
  display: 'block'
};

const Field = React.forwardRef(({placeholder, label, type}, ref) => {
  return (
    <div>
      <label style={labelStyle} >{label}</label>
      <input placeholder={placeholder} ref={ref} type={type} style={inputStyle} />
    </div>
  );
});

const Form = ({onLogin, onRegister}) => {
  const usernameRef = React.useRef();
  const passwordRef = React.useRef();
  var [isLogin, setIsLogin] = React.useState(false);
  const handleSubmit = e => {
      e.preventDefault();
      const data = {
          username: usernameRef.current.value,
          password: passwordRef.current.value
      };
      if (isLogin) {
        onLogin(data);
      } else {
        onRegister(data)
      }
  };
  return (
    <form style={formStyle} onSubmit={handleSubmit} >
      <Field placeholder="Username" ref={usernameRef}  label="Username:" type="text" />
      <Field placeholder="Password" ref={passwordRef} label="Password:" type="password" />
      <div>
        <button onClick={() => (setIsLogin(true))} style={submitStyle} type="primary">{"Login"}</button>
        <button onClick={() => (setIsLogin(false))} style={submitStyle} type="primary">{"Register"}</button>
      </div>
    </form>
  );
};


function Login({updateState, updateToken}) {
    const login_request = (data) => {
      toast("Sending login request to server...");
      fetch(LOGIN_URL,
        {
          method: 'POST',
          headers: { 'content-type': 'application/json' },
          body: JSON.stringify(data)
        }
      )
      .then(response => {
        return response.json();
      })
      .then(function (data) {
        console.log(data)
        toast.dismiss()
        if (data["error"]){
          toast.error(data["error"])
        }else {
          console.log(data["token"])
          toast.success(data["message"])
          updateToken(data["token"])
          updateState(STATES.dashboard)
        }
      })
      .catch(function (error) {
        console.log(error);
      });
    };

    const register_request = (data) => {
      if (data['username'].length == 0) {
        toast.error("Please input valid username")
        return 
      } 
      if (data['password'].length == 0) {
        toast.error("Please input valid password")
        return 
      } 
      if (data['password'].length > 55) {
        toast.error("Username should be less than 56")
        return 
      } 
      if (data['username'].indexOf('$') != -1) {
        toast.error("Username should not have `$`")
        return 
      }
      toast("Sending register request to server...");

      fetch(REGISTER_URL,
        {
          method: 'POST',
          headers: { 'content-type': 'application/json' },
          body: JSON.stringify(data)
        }
      )
      .then(response => {
        return response.json();
      })
      .then(function (data) {
        console.log(data)
        toast.dismiss()
        if (data["error"]){
          toast.error(data["error"])
        }else {
          toast.success(data["message"])
        }
      })
      .catch(function (error) {
        console.log(error);
      });
    };

    return (
      <>
      <div style={appStyle}>
        <Form onRegister={register_request} onLogin={login_request} reg={false}/>
      </div>
      </>
    )
};


export {Login}