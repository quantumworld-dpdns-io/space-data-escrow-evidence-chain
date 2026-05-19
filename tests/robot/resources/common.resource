*** Settings ***
Library    RequestsLibrary
Library    Collections

*** Variables ***
${BASE_URL}    %{BASE_URL=http://localhost:8080}
${API_KEY}     %{API_KEY=dev-api-key}

*** Keywords ***
Create Session
    Create Session    api    ${BASE_URL}

Auth Headers
    [Return]    {"X-API-Key": "${API_KEY}", "Content-Type": "application/json"}
