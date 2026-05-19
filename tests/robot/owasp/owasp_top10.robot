*** Settings ***
Resource    ../resources/common.resource

*** Test Cases ***
A01 Broken Access Control - No API Key
    Init Session
    ${resp}=    GET On Session    api    /v1/search    expected_status=401
    Should Be Equal As Integers    ${resp.status_code}    401

A02 Cryptographic Failures - Verify Missing Resource
    Init Session
    ${headers}=    Auth Headers
    ${resp}=    POST On Session    api    /v1/verify/does-not-exist    headers=${headers}    expected_status=404
    Should Be Equal As Integers    ${resp.status_code}    404

A03 Injection - Search Query Robustness
    Init Session
    ${headers}=    Auth Headers
    ${resp}=    GET On Session    api    /v1/search    headers=    params=q=' OR '1'='1
    Should Be Equal As Integers    ${resp.status_code}    200

A05 Security Misconfiguration - Health Endpoints Exposed but Safe
    Init Session
    ${resp}=    GET On Session    api    /healthz
    Should Be Equal As Integers    ${resp.status_code}    200

A07 Identification and Authentication Failures - Wrong API Key
    Init Session
    ${headers}=    Create Dictionary    X-API-Key=wrong
    ${resp}=    GET On Session    api    /v1/audit    headers=${headers}    expected_status=401
    Should Be Equal As Integers    ${resp.status_code}    401
