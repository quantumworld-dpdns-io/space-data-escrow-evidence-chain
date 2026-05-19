*** Settings ***
Library    RequestsLibrary
Library    Collections

*** Variables ***
${BASE_URL}    %{BASE_URL=http://localhost:8080}
${API_KEY}     %{API_KEY=dev-api-key}

*** Keywords ***
Init Session
    RequestsLibrary.Create Session    api    ${BASE_URL}

Auth Headers
    ${headers}=    Create Dictionary    X-API-Key=${API_KEY}    Content-Type=application/json
    RETURN    ${headers}
