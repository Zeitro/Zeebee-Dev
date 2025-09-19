'use client';
import React, { useMemo, useState, useRef } from 'react';
import PendingPage from './Pending';
const DraftsPage = () => {
	return <PendingPage stage="Draft" key="Draft" />;
};
export default DraftsPage;
