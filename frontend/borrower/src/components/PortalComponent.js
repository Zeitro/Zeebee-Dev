'use client';
import React, { useState, useRef, useEffect } from 'react';
import ReactDOM from 'react-dom';
import eventBus from '@/lib/bus';

// 往根节点portal-root中添加子节点
const PortalComponent = () => {
	const [portalContent, setPortalContent] = useState(null);
	const portalRootRef = useRef(null);

	const ensurePortalRoot = () => {
		if (!portalRootRef.current) {
			const existingRoot = document.getElementById('portal-root');
			if (existingRoot) {
				portalRootRef.current = existingRoot;
			} else {
				portalRootRef.current = document.createElement('div');
				portalRootRef.current.id = 'portal-root';
				document.body.appendChild(portalRootRef.current);
			}
		}
		return portalRootRef.current;
	};

	const onClose = () => {
		setPortalContent(null);
	};

	const appendPortalComponent = (element) => {
		const portalRoot = ensurePortalRoot();
		setPortalContent(ReactDOM.createPortal(element, portalRoot));
	};

	useEffect(() => {
		const handleAppend = (element) => appendPortalComponent(element);
		const handleRemove = () => onClose();

		eventBus?.on('appendPortalComponent', handleAppend);
		eventBus?.on('removePortalComponent', handleRemove);

		return () => {
			eventBus?.removeListener('appendPortalComponent', handleAppend);
			eventBus?.removeListener('removePortalComponent', handleRemove);
			onClose();
			if (portalRootRef.current) {
				document.body.removeChild(portalRootRef.current);
			}
		};
	}, []);

	return portalContent;
};

export default PortalComponent;
