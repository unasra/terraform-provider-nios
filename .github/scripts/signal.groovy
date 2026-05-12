#!/usr/bin/env groovy
/**
 * jenkins/signal.groovy
 * ─────────────────────
 * Drop this file anywhere in your repo (or as a shared library) and load it
 * in your Jenkinsfile with:
 *
 *   def signal = load 'jenkins/signal.groovy'
 *
 * Then call:
 *   signal.success('env-setup')
 *   signal.success('tests-complete', artifacts: ['https://...url1', 'https://...url2'])
 *   signal.error('run-tests', 'go test failed: exit status 1')
 *
 * Required build parameters (passed by the GitHub Action):
 *   SIGNAL_GIST_ID   – Gist ID to write signals into
 *   SIGNAL_TOKEN     – GitHub PAT with gist scope
 *
 * Required build parameter (or env var):
 *   GIT_COMMIT       – the SHA being built (used as a safety key)
 */

// ── private helpers ───────────────────────────────────────────────────────────

/**
 * Push a JSON payload to the signal Gist.
 * Uses a simple sh curl — no plugins required.
 */
private void _push(Map payload) {
    def gistId  = params.SIGNAL_GIST_ID?.trim()
    def token   = params.SIGNAL_TOKEN?.trim()
    def runId   = params.SIGNAL_RUN_ID?.trim() ?: env.SIGNAL_RUN_ID ?: ''
    def sha     = params.GIT_COMMIT?.trim()    ?: env.GIT_COMMIT    ?: ''

    if (!gistId || !token) {
        echo "⚠  SIGNAL_GIST_ID or SIGNAL_TOKEN not set — skipping signal push."
        return
    }

    payload.run_id = runId
    payload.sha    = sha

    def jsonBody = groovy.json.JsonOutput.toJson([
        files: [
            ("ci-signal-${runId}.json".toString()): [
                content: groovy.json.JsonOutput.toJson(payload)
            ]
        ]
    ])

    // Write to a temp file to avoid shell quoting issues with nested JSON
    writeFile file: '.signal_payload.json', text: jsonBody

    def http = sh(
        script: """
            curl -sS -o /dev/null -w "%{http_code}" \\
                -X PATCH \\
                -H "Authorization: token ${token}" \\
                -H "Accept: application/vnd.github+json" \\
                -H "Content-Type: application/json" \\
                --data @.signal_payload.json \\
                "https://api.github.com/gists/${gistId}"
        """,
        returnStdout: true
    ).trim()

    sh "rm -f .signal_payload.json"

    if (http != '200') {
        echo "⚠  Gist PATCH returned HTTP ${http}. The GitHub Action may time out."
    } else {
        echo "📡  Signal pushed → stage='${payload.stage}' status='${payload.status}'"
    }
}

// ── public API ────────────────────────────────────────────────────────────────

/**
 * Send a SUCCESS signal for a stage.
 *
 * @param stageName   matches the STAGE_NAME the Action is waiting for
 * @param opts.message  optional human-readable note
 * @param opts.artifacts  list of direct-download URLs for JUnit XML files
 */
void success(String stageName, Map opts = [:]) {
    _push([
        status   : 'success',
        stage    : stageName,
        message  : opts.message   ?: '',
        artifacts: opts.artifacts ?: []
    ])
}

/**
 * Send an ERROR signal — the GitHub Action will immediately fail with the message.
 *
 * @param stageName   the stage where the failure occurred
 * @param message     human-readable error description
 */
void error(String stageName, String message) {
    _push([
        status   : 'error',
        stage    : stageName,
        message  : message,
        artifacts: []
    ])
}

/**
 * Send an ABORTED signal.
 */
void aborted(String stageName, String message = 'Pipeline was aborted') {
    _push([
        status   : 'aborted',
        stage    : stageName,
        message  : message,
        artifacts: []
    ])
}

// Make functions accessible when the file is loaded with `load`
return this
