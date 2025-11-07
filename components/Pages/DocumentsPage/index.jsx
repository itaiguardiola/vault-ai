import React, { useState, useEffect } from 'react';
import Page from '../../Page';
import Go from '../../Go';

import s from './index.less';

type Props = {
    history: Object,
};

const DocumentsPage = (props: Props) => {
    const [collections, setCollections] = useState([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState('');
    const [deletingUUID, setDeletingUUID] = useState(null);
    const [expandedUUID, setExpandedUUID] = useState(null);
    const [searchTerm, setSearchTerm] = useState('');

    // Load documents on mount
    useEffect(() => {
        loadDocuments();
    }, []);

    const loadDocuments = async () => {
        try {
            setLoading(true);
            setError('');
            const response = await fetch('/api/documents');
            if (!response.ok) {
                throw new Error('Failed to load documents');
            }
            const data = await response.json();
            setCollections(data || []);
        } catch (err) {
            setError(`Error loading documents: ${err.message}`);
        } finally {
            setLoading(false);
        }
    };

    const handleDelete = async (uuid) => {
        if (!window.confirm('Are you sure you want to delete this collection and all its documents? This cannot be undone.')) {
            return;
        }

        try {
            setDeletingUUID(uuid);
            const response = await fetch(`/api/documents/${uuid}`, {
                method: 'DELETE',
            });

            if (!response.ok) {
                throw new Error('Failed to delete collection');
            }

            // Reload documents after successful deletion
            await loadDocuments();
        } catch (err) {
            setError(`Error deleting collection: ${err.message}`);
        } finally {
            setDeletingUUID(null);
        }
    };

    const toggleExpand = (uuid) => {
        setExpandedUUID(expandedUUID === uuid ? null : uuid);
    };

    const filteredCollections = collections.filter((coll) => {
        if (!searchTerm) return true;

        const searchLower = searchTerm.toLowerCase();

        // Search in collection UUID
        if (coll.collection.toLowerCase().includes(searchLower)) return true;

        // Search in document names
        if (coll.documents && coll.documents.some(doc =>
            doc.name.toLowerCase().includes(searchLower)
        )) return true;

        return false;
    });

    const totalDocuments = collections.reduce((sum, coll) =>
        sum + (coll.documents ? coll.documents.length : 0), 0
    );

    const totalChunks = collections.reduce((sum, coll) =>
        sum + (coll.point_count || 0), 0
    );

    if (loading) {
        return (
            <Page history={props.history}>
                <div className={s.documentsPage}>
                    <div className={s.loading}>
                        <div className={s.spinner}></div>
                        <div>Loading documents...</div>
                    </div>
                </div>
            </Page>
        );
    }

    return (
        <Page history={props.history}>
            <div className={s.documentsPage}>
                <div className={s.header}>
                    <div>
                        <h1>📚 Document Library</h1>
                        <p className={s.subtitle}>
                            Manage your uploaded documents and collections
                        </p>
                    </div>
                    <Go to="/" className={s.backButton}>
                        ← Back to Vault
                    </Go>
                </div>

                <div className={s.stats}>
                    <div className={s.statCard}>
                        <div className={s.statValue}>{collections.length}</div>
                        <div className={s.statLabel}>Collections</div>
                    </div>
                    <div className={s.statCard}>
                        <div className={s.statValue}>{totalDocuments}</div>
                        <div className={s.statLabel}>Documents</div>
                    </div>
                    <div className={s.statCard}>
                        <div className={s.statValue}>{totalChunks.toLocaleString()}</div>
                        <div className={s.statLabel}>Total Chunks</div>
                    </div>
                </div>

                <div className={s.controls}>
                    <div className={s.searchBox}>
                        <input
                            type="text"
                            placeholder="🔍 Search documents or collections..."
                            value={searchTerm}
                            onChange={(e) => setSearchTerm(e.target.value)}
                            className={s.searchInput}
                        />
                        {searchTerm && (
                            <button
                                className={s.clearButton}
                                onClick={() => setSearchTerm('')}
                            >
                                ✕
                            </button>
                        )}
                    </div>
                    <button
                        className={s.refreshButton}
                        onClick={loadDocuments}
                        disabled={loading}
                    >
                        ↻ Refresh
                    </button>
                </div>

                {error && (
                    <div className={s.error}>{error}</div>
                )}

                {filteredCollections.length === 0 ? (
                    <div className={s.empty}>
                        <div className={s.emptyIcon}>📭</div>
                        <h2>No documents found</h2>
                        <p>
                            {searchTerm
                                ? `No documents match "${searchTerm}"`
                                : 'Upload some documents to get started!'
                            }
                        </p>
                        <Go to="/" className={s.uploadButton}>
                            Upload Documents
                        </Go>
                    </div>
                ) : (
                    <div className={s.collectionsList}>
                        {filteredCollections.map((collection) => (
                            <div
                                key={collection.collection}
                                className={s.collectionCard}
                            >
                                <div
                                    className={s.collectionHeader}
                                    onClick={() => toggleExpand(collection.collection)}
                                >
                                    <div className={s.collectionInfo}>
                                        <h3>
                                            <span className={s.collectionIcon}>
                                                {expandedUUID === collection.collection ? '📂' : '📁'}
                                            </span>
                                            Collection {collection.collection.substring(0, 8)}...
                                        </h3>
                                        <div className={s.collectionMeta}>
                                            <span className={s.badge}>
                                                {collection.documents ? collection.documents.length : 0} document{collection.documents && collection.documents.length !== 1 ? 's' : ''}
                                            </span>
                                            <span className={s.badge}>
                                                {collection.point_count.toLocaleString()} chunks
                                            </span>
                                        </div>
                                    </div>
                                    <div className={s.collectionActions}>
                                        <button
                                            className={s.deleteButton}
                                            onClick={(e) => {
                                                e.stopPropagation();
                                                handleDelete(collection.collection);
                                            }}
                                            disabled={deletingUUID === collection.collection}
                                        >
                                            {deletingUUID === collection.collection ? '⏳' : '🗑️'} Delete
                                        </button>
                                    </div>
                                </div>

                                {expandedUUID === collection.collection && (
                                    <div className={s.documentsList}>
                                        <h4>Documents in this collection:</h4>
                                        {collection.documents && collection.documents.length > 0 ? (
                                            <ul>
                                                {collection.documents.map((doc, index) => (
                                                    <li key={index} className={s.documentItem}>
                                                        <span className={s.docIcon}>📄</span>
                                                        <span className={s.docName}>{doc.name}</span>
                                                        <span className={s.chunkCount}>
                                                            {doc.chunk_count} chunk{doc.chunk_count !== 1 ? 's' : ''}
                                                        </span>
                                                    </li>
                                                ))}
                                            </ul>
                                        ) : (
                                            <p className={s.noDocuments}>No documents in this collection</p>
                                        )}
                                    </div>
                                )}
                            </div>
                        ))}
                    </div>
                )}

                <div className={s.help}>
                    <h3>💡 Tips</h3>
                    <ul>
                        <li>Each collection is identified by a unique UUID</li>
                        <li>Collections are created automatically when you upload documents</li>
                        <li>Deleting a collection removes all its documents and embeddings</li>
                        <li>Use the search box to find specific documents or collections</li>
                    </ul>
                </div>
            </div>
        </Page>
    );
};

export default DocumentsPage;
