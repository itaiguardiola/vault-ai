import React, { useState, useEffect } from 'react';
import Page from '../../Page';
import Go from '../../Go';

import s from './index.less';

type Props = {
    history: Object,
};

const ConfigPage = (props: Props) => {
    const [config, setConfig] = useState({
        ollama_endpoint: '',
        ollama_embed_model: '',
        ollama_chat_model: '',
        qdrant_endpoint: '',
        max_file_size_mb: 25,
        max_total_size_mb: 50,
    });

    const [loading, setLoading] = useState(true);
    const [saving, setSaving] = useState(false);
    const [testing, setTesting] = useState(false);
    const [message, setMessage] = useState('');
    const [testResults, setTestResults] = useState(null);

    // Load configuration on mount
    useEffect(() => {
        loadConfig();
    }, []);

    const loadConfig = async () => {
        try {
            setLoading(true);
            const response = await fetch('/api/config');
            if (!response.ok) {
                throw new Error('Failed to load configuration');
            }
            const data = await response.json();
            setConfig(data);
            setMessage('');
        } catch (error) {
            setMessage(`Error loading configuration: ${error.message}`);
        } finally {
            setLoading(false);
        }
    };

    const handleInputChange = (field, value) => {
        setConfig({
            ...config,
            [field]: value,
        });
        setMessage('');
    };

    const handleSave = async () => {
        try {
            setSaving(true);
            setMessage('');

            const response = await fetch('/api/config', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(config),
            });

            if (!response.ok) {
                throw new Error('Failed to save configuration');
            }

            const data = await response.json();
            setMessage(data.message || 'Configuration saved successfully!');
        } catch (error) {
            setMessage(`Error saving configuration: ${error.message}`);
        } finally {
            setSaving(false);
        }
    };

    const handleTest = async () => {
        try {
            setTesting(true);
            setTestResults(null);

            const response = await fetch('/api/config/test');
            if (!response.ok) {
                throw new Error('Failed to test connections');
            }

            const data = await response.json();
            setTestResults(data);
        } catch (error) {
            setMessage(`Error testing connections: ${error.message}`);
        } finally {
            setTesting(false);
        }
    };

    const handleReset = () => {
        loadConfig();
        setMessage('Configuration reset to saved values');
        setTestResults(null);
    };

    if (loading) {
        return (
            <Page history={props.history}>
                <div className={s.configPage}>
                    <div className={s.loading}>Loading configuration...</div>
                </div>
            </Page>
        );
    }

    return (
        <Page history={props.history}>
            <div className={s.configPage}>
                <div className={s.header}>
                    <h1>⚙️ Configuration</h1>
                    <p className={s.subtitle}>
                        Configure your local Askara installation
                    </p>
                </div>

                <div className={s.configForm}>
                    <div className={s.section}>
                        <h2>🤖 Ollama Configuration</h2>
                        <p className={s.sectionDesc}>
                            Local LLM service for embeddings and chat
                        </p>

                        <div className={s.formGroup}>
                            <label>Ollama Endpoint</label>
                            <input
                                type="text"
                                value={config.ollama_endpoint}
                                onChange={(e) =>
                                    handleInputChange('ollama_endpoint', e.target.value)
                                }
                                placeholder="http://ollama:11434"
                            />
                            <span className={s.hint}>
                                Default: http://ollama:11434 (Docker) or http://localhost:11434 (local)
                            </span>
                        </div>

                        <div className={s.formGroup}>
                            <label>Embedding Model</label>
                            <input
                                type="text"
                                value={config.ollama_embed_model}
                                onChange={(e) =>
                                    handleInputChange('ollama_embed_model', e.target.value)
                                }
                                placeholder="nomic-embed-text"
                            />
                            <span className={s.hint}>
                                Model for generating document embeddings (768 dimensions)
                            </span>
                        </div>

                        <div className={s.formGroup}>
                            <label>Chat Model</label>
                            <input
                                type="text"
                                value={config.ollama_chat_model}
                                onChange={(e) =>
                                    handleInputChange('ollama_chat_model', e.target.value)
                                }
                                placeholder="llama3"
                            />
                            <span className={s.hint}>
                                Model for generating answers (llama3, mistral, etc.)
                            </span>
                        </div>
                    </div>

                    <div className={s.section}>
                        <h2>🗄️ Qdrant Configuration</h2>
                        <p className={s.sectionDesc}>
                            Local vector database for storing embeddings
                        </p>

                        <div className={s.formGroup}>
                            <label>Qdrant Endpoint</label>
                            <input
                                type="text"
                                value={config.qdrant_endpoint}
                                onChange={(e) =>
                                    handleInputChange('qdrant_endpoint', e.target.value)
                                }
                                placeholder="http://qdrant:6333"
                            />
                            <span className={s.hint}>
                                Default: http://qdrant:6333 (Docker) or http://localhost:6333 (local)
                            </span>
                        </div>
                    </div>

                    <div className={s.section}>
                        <h2>📁 File Upload Limits</h2>
                        <p className={s.sectionDesc}>
                            Maximum file sizes for uploads
                        </p>

                        <div className={s.formGroup}>
                            <label>Max File Size (MB)</label>
                            <input
                                type="number"
                                value={config.max_file_size_mb}
                                onChange={(e) =>
                                    handleInputChange('max_file_size_mb', parseInt(e.target.value))
                                }
                                min="1"
                                max="1000"
                            />
                            <span className={s.hint}>
                                Maximum size for individual files
                            </span>
                        </div>

                        <div className={s.formGroup}>
                            <label>Max Total Upload Size (MB)</label>
                            <input
                                type="number"
                                value={config.max_total_size_mb}
                                onChange={(e) =>
                                    handleInputChange('max_total_size_mb', parseInt(e.target.value))
                                }
                                min="1"
                                max="5000"
                            />
                            <span className={s.hint}>
                                Maximum total size for batch uploads
                            </span>
                        </div>
                    </div>

                    {message && (
                        <div className={message.includes('Error') ? s.errorMessage : s.successMessage}>
                            {message}
                        </div>
                    )}

                    {testResults && (
                        <div className={s.testResults}>
                            <h3>Connection Test Results</h3>
                            <div className={s.testResult}>
                                <span className={s.testLabel}>Ollama:</span>
                                <span className={testResults.ollama.status === 'ok' ? s.statusOk : s.statusError}>
                                    {testResults.ollama.status === 'ok' ? '✓' : '✗'} {testResults.ollama.message}
                                </span>
                            </div>
                            <div className={s.testResult}>
                                <span className={s.testLabel}>Qdrant:</span>
                                <span className={testResults.qdrant.status === 'ok' ? s.statusOk : s.statusError}>
                                    {testResults.qdrant.status === 'ok' ? '✓' : '✗'} {testResults.qdrant.message}
                                </span>
                            </div>
                        </div>
                    )}

                    <div className={s.actions}>
                        <button
                            className={s.testButton}
                            onClick={handleTest}
                            disabled={testing}
                        >
                            {testing ? 'Testing...' : '🔍 Test Connections'}
                        </button>
                        <button
                            className={s.resetButton}
                            onClick={handleReset}
                            disabled={saving}
                        >
                            ↺ Reset
                        </button>
                        <button
                            className={s.saveButton}
                            onClick={handleSave}
                            disabled={saving}
                        >
                            {saving ? 'Saving...' : '💾 Save Configuration'}
                        </button>
                    </div>

                    <div className={s.helpSection}>
                        <h3>ℹ️ Help</h3>
                        <ul>
                            <li>
                                <strong>Docker Setup:</strong> Use service names (ollama, qdrant) for endpoints
                            </li>
                            <li>
                                <strong>Local Setup:</strong> Use localhost:PORT for endpoints
                            </li>
                            <li>
                                <strong>Models:</strong> Make sure models are pulled in Ollama before use
                            </li>
                            <li>
                                <strong>Restart Required:</strong> Changes take effect after restarting the application
                            </li>
                        </ul>
                    </div>
                </div>

                <Go to="/" className={s.backLink}>
                    ← Back to Vault
                </Go>
            </div>
        </Page>
    );
};

export default ConfigPage;
