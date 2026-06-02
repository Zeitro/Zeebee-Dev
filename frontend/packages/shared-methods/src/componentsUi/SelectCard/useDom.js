import React from 'react';
import PropertyCardItem from './Item/property';
import BorrowerCardItem from './Item/borrower';
import BorrowingEntityCardItem from './Item/borrowingEntity';
import TransactionCardItem from './Item/transaction';
import PropertyAdd from './Add/property';
import BorrowerAdd from './Add/borrower';
import BorrowingEntityAdd from './Add/borrowingEntity';
const useDom = () => {
	const getCardItem = (type, config = {}) => {
		switch (type) {
			case 'property':
				return <PropertyCardItem key={config.index} {...config} />;
			case 'borrower':
				return <BorrowerCardItem key={config.index} {...config} />;
			case 'borrowingEntity':
				return <BorrowingEntityCardItem key={config.index} {...config} />;
			case 'transaction':
				return <TransactionCardItem key={config.index} {...config} />;

			default:
				break;
		}
	};

	const getAddDom = (type, config = {}) => {
		switch (type) {
			case 'property':
				return <PropertyAdd key={config.key} {...config} />;
			case 'borrower':
				return <BorrowerAdd key={config.key} {...config} />;
			case 'borrowingEntity':
				return <BorrowingEntityAdd key={config.key} {...config} />;
			default:
				break;
		}
	};
	return {
		getAddDom,
		getCardItem,
	};
};
export default useDom;
