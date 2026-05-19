*** Settings ***
Resource    ../resources/common.resource

*** Test Cases ***
A01 Broken Access Control - No API Key
    Init Session
    ${resp}=    GET On Session    api    /v1/search    expected_status=401
    Should Be Equal As Integers    ${resp.status_code}    401

A01 Broken Access Control - Viewer Cannot Access Admin
    Init Session
    ${headers}=    Create Dictionary    Authorization=Bearer jwt-viewer
    ${resp}=    GET On Session    api    /v1/admin/key-rotation    headers=${headers}    expected_status=403
    Should Be Equal As Integers    ${resp.status_code}    403

A02 Cryptographic Failures - Verify Missing Resource
    Init Session
    ${headers}=    Auth Headers
    ${resp}=    POST On Session    api    /v1/verify/does-not-exist    headers=${headers}    expected_status=404
    Should Be Equal As Integers    ${resp.status_code}    404

A02 Cryptographic Failures - Hash Material Present
    Init Session
    ${headers}=    Auth Headers
    ${body}=    Set Variable    {"external_id":"EXT-CRYPTO","source":"sat-c","type":"imagery","payload":{"file":"f.tif"}}
    ${created}=    POST On Session    api    /v1/evidence    data=${body}    headers=${headers}
    Should Be True    '${created.json()}[hash]' != ''

A03 Injection - Search Query Robustness
    Init Session
    ${headers}=    Auth Headers
    ${resp}=    GET On Session    api    /v1/search    headers=${headers}    params=q=' OR '1'='1
    Should Be Equal As Integers    ${resp.status_code}    200

A04 Insecure Design - Reject Invalid Attestation
    Init Session
    ${headers}=    Auth Headers
    ${body}=    Set Variable    {"evidence_id":"","signer":"","signature":"","algorithm":""}
    ${resp}=    POST On Session    api    /v1/attest    data=${body}    headers=${headers}    expected_status=400
    Should Be Equal As Integers    ${resp.status_code}    400

A05 Security Misconfiguration - Health Endpoints Exposed but Safe
    Init Session
    ${resp}=    GET On Session    api    /healthz
    Should Be Equal As Integers    ${resp.status_code}    200

A06 Vulnerable Components - Version Endpoint Present
    Init Session
    ${resp}=    GET On Session    api    /version
    Should Be Equal As Integers    ${resp.status_code}    200

A07 Identification and Authentication Failures - Wrong API Key
    Init Session
    ${headers}=    Create Dictionary    X-API-Key=wrong
    ${resp}=    GET On Session    api    /v1/audit    headers=${headers}    expected_status=401
    Should Be Equal As Integers    ${resp.status_code}    401

A08 Software and Data Integrity Failures - Idempotency Replay
    Init Session
    ${headers}=    Auth Headers
    Set To Dictionary    ${headers}    Idempotency-Key=owasp-a08-1
    ${body}=    Set Variable    {"external_id":"EXT-IDEM-A08","source":"sat-i","type":"imagery","payload":{"k":"v"}}
    ${r1}=    POST On Session    api    /v1/evidence    data=${body}    headers=${headers}
    ${r2}=    POST On Session    api    /v1/evidence    data=${body}    headers=${headers}
    Should Be Equal As Integers    ${r1.status_code}    201
    Should Be Equal As Integers    ${r2.status_code}    200

A09 Security Logging and Monitoring - Audit Endpoint Reachable
    Init Session
    ${headers}=    Auth Headers
    ${resp}=    GET On Session    api    /v1/audit    headers=${headers}
    Should Be Equal As Integers    ${resp.status_code}    200

A10 SSRF - External URL Query String Does Not Break API
    Init Session
    ${headers}=    Auth Headers
    ${resp}=    GET On Session    api    /v1/search    headers=${headers}    params=q=http://169.254.169.254/latest/meta-data
    Should Be Equal As Integers    ${resp.status_code}    200
