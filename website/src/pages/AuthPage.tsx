import { useState, useEffect } from 'react';


interface IAuth {
  authorised: boolean;
  name: string;
  email: string;
  company: string;
}

// const CLIEndpoint = "http://localhost:21765";
const CLIEndpoint = "https://jsonplaceholder.typicode.com/users/1";

const callAuthEndpoint = async (code: string | null, state: string | null) => {
  // do this in the go app just an example here!
  const response = await fetch(CLIEndpoint)
  const data = await response.json()
  console.log("data: ")
  console.log(data)

  const authResponse: IAuth = {
    authorised: true,
    name: data.name,
    email: data.email,
    company: data.company.name,
  }
  return authResponse;
}

const AuthPage = () => {
  const [auth, setAuth] = useState<IAuth>();
  const params = new URLSearchParams(window.location.search);
  const code = params.get('code');
  const state = params.get('state');
  const localState = localStorage.getItem('gitHubState');
  const getDetailsList = () => {
    if (!auth) return
    return Object.entries(auth as IAuth).map(([key, val], i) => {
      if (key == 'authorised') return
      return <li key={i}><span>{key}: {val}</span></li>
    })
  }

  useEffect(() => {
    if (!auth?.name) {
      console.log("calling auth endpoint")
      callAuthEndpoint(code, state)
        .then(res => {
          setAuth(res)
        })
    }
  })

  //TODO: If local states don't match throw auth error
  console.log("Matches: " + (localState == state))

  return (
    <div>
      <header className="App-header">
        {auth?.authorised &&
          <div>
            <h2>
              Hi {auth.name}, you are now Authorised!
            </h2>
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