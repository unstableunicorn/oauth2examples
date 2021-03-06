import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import React from 'react';
import { faGithub } from '@fortawesome/free-brands-svg-icons';
import styled from 'styled-components';
import {IAuthParams} from '../pages/HomePage'

const StyledButton = styled.button`
  font-size: medium;
  padding-left: 10px;
  color: #fff;
  background-color: #444;
  border-color: rgba(0,0,0,0.2);
  position: relative;
  text-align: left;
  white-space: nowrap;
  display: block;
`;

const HoverWrapper = styled.div`
  &:hover ${StyledButton} {
    background-color: #2b2b2b;
  }
`;

const StyledButtonText = styled.span`
  font-size: medium;
  text-align: center;
  padding-left: 10px;
  border-left: 1px solid rgba(0,0,0,0.2);
`;

const StyledFontAwesomeIcon = styled.i`
  padding-right: 10px;
`;

const GitHubLoginURL = "https://github.com/login/oauth/authorize"

const getGitHubParams = () => {
  const authParams: IAuthParams = JSON.parse(localStorage.getItem('authParams') as string)
  return new URLSearchParams({
    client_id: authParams.clientId,
    redirect_uri: 'http://localhost:3000/auth',
    scope: authParams.scope,
    state: authParams.state
  });
}

const onClick = () => {
  window.location.replace(`${GitHubLoginURL}?${getGitHubParams()}`);
}

const GitHubButton = () => {
  return (
    <React.Fragment>
      <HoverWrapper>
        <StyledButton
          className="title active"
          onClick={() => onClick()}
        >
          <StyledFontAwesomeIcon>
            <FontAwesomeIcon className="fa-lg" icon={faGithub} />
          </StyledFontAwesomeIcon>
          <StyledButtonText>Sign in to GitHub</StyledButtonText>
        </StyledButton>
      </HoverWrapper>
    </React.Fragment>
  )
}

export default GitHubButton