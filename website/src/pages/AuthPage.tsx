import { useState, useEffect } from 'react';
import styled from 'styled-components';
import { IAuthParams } from './HomePage';


interface IAuth {
  authorised: boolean;
  name: string;
  alias: string;
  email: string;
  avatar: string;
}

const AvatarStyle = styled.img`
  position: relative;
  width: 80px;
  float: left ;
`
const Header2Style = styled.h2`
  float: right;
`

const CLIEndpoint = "http://localhost:3001/oauth/callback";

const callAuthEndpoint = async (code: string | null, state: string | null) => {
  // do this in the go app just an example here!
  const response = await fetch(`${CLIEndpoint}?code=${code}`)
  console.log("awaiting response")
  const data = await response.json()
  console.log("data: ")
  console.log(data)

  const authResponse: IAuth = {
    authorised: true,
    name: data.name,
    email: data.email,
    alias: data.login,
    avatar: data.avatar_url
  }
  return authResponse;
}

const AuthPage = () => {
  const [auth, setAuth] = useState<IAuth>();
  const params = new URLSearchParams(window.location.search);
  const code = params.get('code');
  const state = params.get('state');
  const localState: IAuthParams = JSON.parse(localStorage.getItem('authParams') as string);
  const getDetailsList = () => {
    if (!auth) return
    return Object.entries(auth as IAuth).map(([key, val], i) => {
      if (key == 'authorised') return
      return <li key={i}><span>{key}: {val}</span></li>
    })
  }
  // just a thing
  useEffect(() => {
    if (!auth?.name) {
      console.log("calling auth endpoint")
      callAuthEndpoint(code, state)
        .then(res => {
          console.log("Got reponse:")
          console.log(res)
          setAuth(res)
        })
    }
  })

  //TODO: If local states don't match throw auth error
  console.log("Matches: " + (localState.state == state))

  return (
    <div>
      <header className="App-header">
        {auth?.authorised &&
          <div>
            <div className="header">
            <AvatarStyle src={auth.avatar} alt={auth.alias}/>
            <Header2Style>
              Hi {auth.name}, You are now Authorised!
            </Header2Style>
            </div>
              <h3>User Details: </h3>
              <ul>
                {getDetailsList()}
              </ul>
            <p>
              You can now close this page!
            </p>
          </div>
        }
        {!auth?.authorised && <p>
          Waiting on Authorisation...
        </p>}
      </header>
    </div >
  )
}

export default AuthPage;