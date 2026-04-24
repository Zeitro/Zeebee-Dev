'use client';

import React, { useCallback, useEffect, useMemo, useState } from 'react';
import { Button, Modal, Space, Spin, Table, Tag, Typography } from 'antd';
import Cookie from 'js-cookie';

const { Text } = Typography;

const DISCOVER_TIMEOUT_MS = 5 * 60 * 1000;

async function longRequest(path, { method = 'POST', body } = {}) {
	const base = process.env.NEXT_PUBLIC_API_URL || '';
	const url = `${base}/api/v1${path}`;
	const store = JSON.parse(
		Cookie.get(`permissionStore_${process.env.NEXT_PUBLIC_SYSTEM_TYPE}`) || '{}',
	);
	const token = store?.state?.access?.access_token;
	const res = await fetch(url, {
		method,
		headers: {
			'Content-Type': 'application/json',
			...(token ? { Authorization: `Bearer ${token}` } : {}),
		},
		body: body ? JSON.stringify(body) : undefined,
		signal: AbortSignal.timeout(DISCOVER_TIMEOUT_MS),
	});
	const json = await res.json();
	if (json?.code === 401) {
		window.location.href = '/login';
		throw new Error('401');
	}
	if (json?.code && json.code !== 200) throw new Error(json.message || `fail:${json.code}`);
	return json.data;
}

function fmtUsdRange(r) {
	if (!r || (r.min === 0 && r.max === 0)) return '$0 - $0';
	return `$${r.min} - $${r.max}`;
}

function fmtMonthsRange(r) {
	if (!r || (r.min === 0 && r.max === 0)) return '0 - 0 months';
	return `${r.min} - ${r.max} months`;
}

function fmtPctRange(r) {
	if (!r || (r.min === 0 && r.max === 0)) return '0 - 0%';
	return `${r.min} - ${r.max}%`;
}

function candidateToCommitLender(c) {
	const lm = c.lender_metadata || {};
	return {
		name: c.name,
		contact_email: c.contact_email,
		contact_phone: c.contact_phone,
		contact_name: c.contact_name,
		website: c.website,
		tagline: c.tagline,
		logo: c.logo,
		address: c.address || {},
		lender_metadata: {
			loan_amount_range: lm.loan_amount_range || { min: 0, max: 0 },
			loan_term_range: lm.loan_term_range || { min: 0, max: 0 },
			loan_rate_range: lm.loan_rate_range || { min: 0, max: 0 },
			loan_point_range: lm.loan_point_range || { min: 0, max: 0 },
		},
	};
}

/**
 * Manual lender discover: preview OpenAI results, then commit selected "new" rows only.
 */
export default function DiscoverPreviewModal({ open, onClose, onCommitted }) {
	const [loading, setLoading] = useState(false);
	const [committing, setCommitting] = useState(false);
	const [candidates, setCandidates] = useState([]);
	const [selectedKeys, setSelectedKeys] = useState([]);

	const loadPreview = useCallback(async () => {
		setLoading(true);
		setCandidates([]);
		setSelectedKeys([]);
		try {
			const res = await longRequest('/lender_discover/preview', { method: 'POST', body: {} });
			const list = res?.data?.candidates ?? [];
			setCandidates(list);
			const defaults = list
				.map((c, i) => (c.dedupe_status === 'new' ? String(i) : null))
				.filter(Boolean);
			setSelectedKeys(defaults);
		} catch (e) {
			console.error(e);
		} finally {
			setLoading(false);
		}
	}, []);

	useEffect(() => {
		if (!open) return;
		loadPreview();
	}, [open, loadPreview]);

	const dataSource = useMemo(
		() => candidates.map((c, i) => ({ ...c, key: String(i) })),
		[candidates],
	);

	const columns = useMemo(
		() => [
			{
				title: 'Lender',
				key: 'name',
				width: 220,
				render: (_, row) => (
					<Space size={8} wrap>
						<Text strong={row.dedupe_status === 'new'}>{row.name || '—'}</Text>
						{row.dedupe_status === 'new' ? (
							<Tag color="green">New</Tag>
						) : (
							<Tag>Exists</Tag>
						)}
					</Space>
				),
			},
			{
				title: 'Loan amount (US $)',
				key: 'amt',
				width: 140,
				render: (_, row) => fmtUsdRange(row.lender_metadata?.loan_amount_range),
			},
			{
				title: 'Loan term (months)',
				key: 'term',
				width: 150,
				render: (_, row) => fmtMonthsRange(row.lender_metadata?.loan_term_range),
			},
			{
				title: 'Loan rate',
				key: 'rate',
				width: 110,
				render: (_, row) => fmtPctRange(row.lender_metadata?.loan_rate_range),
			},
			{
				title: 'Loan points',
				key: 'pts',
				width: 110,
				render: (_, row) => fmtPctRange(row.lender_metadata?.loan_point_range),
			},
			{
				title: 'Website',
				dataIndex: 'website',
				key: 'website',
				ellipsis: true,
				render: (v) => v || '—',
			},
			{
				title: 'Contact',
				key: 'contact',
				width: 180,
				ellipsis: true,
				render: (_, row) => {
					const parts = [row.contact_email, row.contact_phone].filter(Boolean);
					return parts.length ? parts.join(' · ') : '—';
				},
			},
		],
		[],
	);

	const handleCommit = async () => {
		const lenders = selectedKeys
			.map((k) => candidates[Number(k)])
			.filter(Boolean)
			.map(candidateToCommitLender);
		if (!lenders.length) {
			onClose?.();
			return;
		}
		setCommitting(true);
		try {
			await longRequest('/lender_discover/commit', {
				method: 'POST',
				body: { lenders },
			});
			onCommitted?.();
			onClose?.();
		} catch (e) {
			console.error(e);
		} finally {
			setCommitting(false);
		}
	};

	return (
		<Modal
			title="Search latest lenders from ChatGPT"
			open={open}
			onCancel={onClose}
			width={960}
			styles={{ body: { paddingTop: 12 } }}
			footer={
				<Space>
					<Button onClick={onClose}>Cancel</Button>
					<Button
						type="primary"
						loading={committing}
						disabled={loading}
						onClick={handleCommit}
					>
						Update
					</Button>
				</Space>
			}
			destroyOnClose
		>
			{loading ? (
				<div className="flex flex-col items-center justify-center gap-3 py-16">
					<Spin size="large" />
					<Text type="secondary">Searching ChatGPT…</Text>
				</div>
			) : (
				<Table
					size="small"
					pagination={false}
					scroll={{ x: 900, y: 420 }}
					rowClassName={(record) =>
						record.dedupe_status === 'existed' ? 'opacity-50' : ''
					}
					rowSelection={{
						selectedRowKeys: selectedKeys,
						onChange: setSelectedKeys,
						getCheckboxProps: (record) => ({
							disabled: record.dedupe_status !== 'new',
						}),
					}}
					columns={columns}
					dataSource={dataSource}
				/>
			)}
		</Modal>
	);
}
