# Implementation Plan - 300 One-Commit TODOs

This plan defines exactly 300 TODO items, each intended to be delivered in one atomic commit.

1. Create `docs/IMPLEMENTATION_PLAN_300_COMMITS.md`.
2. Add architecture overview diagram source.
3. Add ADR template.
4. Add ADR-001 Go+Gin choice.
5. Add ADR-002 DuckDB+Qdrant choice.
6. Add ADR-003 Ollama integration decision.
7. Add ADR-004 OTel+Weave decision.
8. Initialize Go module.
9. Add Makefile baseline.
10. Add editor and formatting rules.
11. Add API bootstrap entrypoint.
12. Add config loader package.
13. Add env default policy.
14. Add structured logging.
15. Add `/healthz` endpoint.
16. Add `/readyz` endpoint.
17. Add graceful shutdown flow.
18. Add version metadata endpoint.
19. Add `/v1` route grouping.
20. Add middleware chain.
21. Add request-id middleware.
22. Add recovery middleware.
23. Add timeout middleware.
24. Add CORS middleware.
25. Add API error envelope.
26. Add validation error mapper.
27. Add domain package skeleton.
28. Add `EvidenceRecord` type.
29. Add `CustodyEvent` type.
30. Add `Attestation` type.
31. Add `VerificationReport` type.
32. Add hash interface.
33. Add signer interface.
34. Add repository interfaces.
35. Add service skeleton.
36. Add in-memory repositories.
37. Add evidence create logic.
38. Add evidence get logic.
39. Add custody append logic.
40. Add chain verify logic.
41. Add POST `/v1/evidence` handler.
42. Add GET `/v1/evidence/:id` handler.
43. Add POST `/v1/custody` handler.
44. Add POST `/v1/verify/:id` handler.
45. Add evidence request DTOs.
46. Add evidence response DTOs.
47. Add custody DTOs.
48. Add verify DTOs.
49. Add router tests.
50. Add evidence handler tests.
51. Add custody handler tests.
52. Add verify handler tests.
53. Add service unit tests.
54. Add custody service tests.
55. Add verify service tests.
56. Add repository contract tests.
57. Add OpenAPI skeleton.
58. Add OpenAPI health paths.
59. Add OpenAPI evidence paths.
60. Add OpenAPI custody paths.
61. Add OpenAPI verify paths.
62. Add OpenAPI error schemas.
63. Add OpenAPI core schemas.
64. Add OpenAPI examples.
65. Add OpenAPI lint config.
66. Add OpenAPI CI validation.
67. Add generated stubs check.
68. Add generated Go client.
69. Add API conformance tests.
70. Add crypto package skeleton.
71. Add SHA-256 implementation.
72. Add SHA-3 implementation.
73. Add BLAKE3 implementation.
74. Add hash selection config.
75. Add ECDSA signer.
76. Add Ed25519 signer.
77. Add key loading utility.
78. Add signature envelope.
79. Add signature verify utility.
80. Add crypto unit tests.
81. Add sign/verify tests.
82. Add deterministic fixtures.
83. Add chain-link model.
84. Add payload canonicalizer.
85. Add canonicalization tests.
86. Add immutable chain builder.
87. Add integrity checker.
88. Add tamper-detection tests.
89. Add timestamp normalizer.
90. Add UTC RFC3339 policy.
91. Add duplicate detection policy.
92. Add idempotency support.
93. Add idempotency middleware.
94. Add idempotency tests.
95. Add pagination support.
96. Add source/type/date filtering.
97. Add sorting support.
98. Add query parser.
99. Add query parser tests.
100. Add list endpoint tests.
101. Add DuckDB adapter skeleton.
102. Add DuckDB connection manager.
103. Add migration framework.
104. Add migration evidence table.
105. Add migration custody table.
106. Add migration attestation table.
107. Add migration audit table.
108. Add migration runner command.
109. Add DuckDB evidence repo.
110. Add DuckDB custody repo.
111. Add DuckDB attestation repo.
112. Add DuckDB repo tests.
113. Add DB transactions helper.
114. Add rollback tests.
115. Add Qdrant adapter skeleton.
116. Add Qdrant collection bootstrap.
117. Add embedding mapping.
118. Add vector upsert flow.
119. Add semantic search flow.
120. Add metadata filter mapping.
121. Add hybrid retrieval hook.
122. Add Qdrant integration tests.
123. Add index-on-create.
124. Add index-on-enrich.
125. Add index-delete-on-purge.
126. Add GET `/v1/search` endpoint.
127. Add search ranking metadata.
128. Add search validation.
129. Add search handler tests.
130. Add search service tests.
131. Add search integration tests.
132. Add Ollama client skeleton.
133. Add Ollama health probe.
134. Add Ollama generate wrapper.
135. Add prompt template module.
136. Add classify prompt template.
137. Add summarize prompt template.
138. Add extract prompt template.
139. Add model endpoint config.
140. Add AI feature flags.
141. Add enrichment orchestration.
142. Add enrichment state model.
143. Add enrichment trigger endpoint.
144. Add enrichment status endpoint.
145. Add enrichment unit tests.
146. Add enrichment integration tests.
147. Add Ollama unavailable fallback.
148. Add retry policy.
149. Add AI endpoint rate limit.
150. Add prompt redaction utility.
151. Add auth package skeleton.
152. Add API-key middleware.
153. Add JWT middleware.
154. Add role model.
155. Add RBAC evaluator.
156. Add route authorization wiring.
157. Add auth failure audit log.
158. Add auth unit tests.
159. Add RBAC unit tests.
160. Add key rotation metadata endpoint.
161. Add chain timeline endpoint.
162. Add attestation endpoint.
163. Add verify endpoint improvements.
164. Add bulk verify endpoint.
165. Add verify serialization model.
166. Add proof bundle export.
167. Add proof bundle docs.
168. Add export tests.
169. Add retention policy model.
170. Add retention evaluator.
171. Add purge scheduler.
172. Add purge dry-run.
173. Add purge audit logging.
174. Add purge tests.
175. Add audit query endpoint.
176. Add audit filtering.
177. Add audit cursor pagination.
178. Add audit CSV export.
179. Add audit JSONL export.
180. Add audit tests.
181. Add OTel tracing setup.
182. Add OTel metrics setup.
183. Add log-trace correlation.
184. Add HTTP tracing middleware.
185. Add DuckDB spans.
186. Add Qdrant spans.
187. Add Ollama spans.
188. Add evidence lifecycle metrics.
189. Add verification metrics.
190. Add security metrics.
191. Add telemetry unit tests.
192. Add telemetry smoke tests.
193. Add Weave adapter skeleton.
194. Add Weave trace emission.
195. Add eval dataset schema.
196. Add evaluation scaffold.
197. Add extraction evaluator.
198. Add summary evaluator.
199. Add nightly eval command.
200. Add eval artifact report.
201. Refactor CI stages.
202. Add lint stage.
203. Add coverage stage.
204. Add integration matrix stage.
205. Add OpenAPI validation stage.
206. Add dependency scan stage.
207. Add secret scan stage.
208. Add Go SAST stage.
209. Add license check stage.
210. Add artifact build stage.
211. Add image build stage.
212. Add SBOM generation stage.
213. Add provenance stage.
214. Add release tagging workflow.
215. Add changelog workflow.
216. Add PR quality gates workflow.
217. Add branch protection docs.
218. Add env promotion workflow.
219. Add staging deploy workflow.
220. Add production deploy workflow.
221. Add rollback workflow.
222. Add Robot framework skeleton.
223. Add Robot resource keywords.
224. Add Robot auth keyword.
225. Add Robot evidence keyword.
226. Add Robot custody keyword.
227. Add Robot verify keyword.
228. Add Robot happy-path suite.
229. Add Robot validation-negative suite.
230. Add Robot authz suite.
231. Add Robot idempotency suite.
232. Add Robot pagination suite.
233. Add Robot enrichment suite.
234. Add Robot audit/export suite.
235. Add Robot CI job.
236. Add Robot reports artifact.
237. Add security harness structure.
238. Add OWASP A01 tests.
239. Add OWASP A02 tests.
240. Add OWASP A03 tests.
241. Add OWASP A04 tests.
242. Add OWASP A05 tests.
243. Add OWASP A06 tests.
244. Add OWASP A07 tests.
245. Add OWASP A08 tests.
246. Add OWASP A09 tests.
247. Add OWASP A10 tests.
248. Add DAST-lite script.
249. Add ZAP baseline CI job.
250. Add security gating threshold.
251. Add Tetragon profile docs.
252. Add suspicious exec policy sample.
253. Add egress anomaly policy sample.
254. Add runtime event mapping.
255. Add runtime security tests.
256. Add PQC interface.
257. Add Dilithium placeholder adapter.
258. Add Kyber placeholder adapter.
259. Add PQC config negotiation.
260. Add dual-sign mode.
261. Add PQC compatibility tests.
262. Add MCP manifest.
263. Add MCP adapter skeleton.
264. Add MCP evidence ingest action.
265. Add MCP verification action.
266. Add MCP search action.
267. Add MCP audit action.
268. Add OpenAPI tool-calling notes.
269. Add MCP integration tests.
270. Add CLI skeleton.
271. Add CLI evidence-create.
272. Add CLI evidence-verify.
273. Add CLI audit-query.
274. Add CLI enrich-trigger.
275. Add CLI tests.
276. Add API Dockerfile.
277. Add docker-compose full stack.
278. Add local bootstrap script.
279. Add `make dev-up`.
280. Add `make dev-down`.
281. Add `make test-all`.
282. Add `make security-test`.
283. Add local development docs.
284. Add API usage docs.
285. Add integration architecture docs.
286. Add operations runbook.
287. Add incident response runbook.
288. Add key management runbook.
289. Add backup/restore runbook.
290. Add retention/compliance docs.
291. Add threat model docs.
292. Add SLO/alerts docs.
293. Add load test skeleton.
294. Add performance baseline job.
295. Add regression gate.
296. Add final end-to-end suite.
297. Add release candidate checklist.
298. Add go-live checklist.
299. Add post-deploy verification.
300. Add governance sign-off artifact.
