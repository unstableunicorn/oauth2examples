import GitHubButton from '../components/githubbutton'

type NullString = string | null;

export interface IAuthParams {
  clientId: string,
  state: string,
  scope: string,
  redirectUri: string,
}
const HomePage = () => {
  const params = new URLSearchParams(window.location.search);
  const clientId = params.get('client_id')
  const state = params.get('state')
  const scope = params.get('scope')
  const redirectUri = params.get('redirect_uri')

  const authParams: IAuthParams = {
    clientId: clientId ? clientId : '8f300287d12718a77080',
    state: state ? state : '',
    scope: scope ? scope : 'read:user read:email',
    redirectUri: redirectUri ? redirectUri : '',
  }

  localStorage.setItem('authParams', JSON.stringify(authParams))
  return (
    <div className="App">
      <header className="App-header">
        <p>
          Log in to the App
        </p>
        <GitHubButton />
      </header>
    </div>
  )
}

export default HomePage;