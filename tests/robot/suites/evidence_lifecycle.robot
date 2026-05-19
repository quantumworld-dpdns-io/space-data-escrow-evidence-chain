*** Settings ***
Resource    ../resources/common.resource

*** Test Cases ***
Evidence Lifecycle Happy Path
    Create Session
    ${headers}=    Auth Headers
    ${create_body}=    Set Variable    {"external_id":"EXT-R1","source":"sat-r","type":"imagery","payload":{"file":"r1.tif"}}
    ${create_resp}=    POST On Session    api    /v1/evidence    data=${create_body}    headers=${headers}
    Should Be Equal As Integers    ${create_resp.status_code}    201
    ${id}=    Set Variable    ${create_resp.json()}[id]

    ${custody_body}=    Set Variable    {"evidence_id":"${id}","actor":"robot","action":"ingest","note":"robot path"}
    ${custody_resp}=    POST On Session    api    /v1/custody    data=${custody_body}    headers=${headers}
    Should Be Equal As Integers    ${custody_resp.status_code}    202

    ${verify_resp}=    POST On Session    api    /v1/verify/${id}    headers=${headers}
    Should Be Equal As Integers    ${verify_resp.status_code}    200
    Should Be Equal    ${verify_resp.json()}[chain_valid]    ${True}
