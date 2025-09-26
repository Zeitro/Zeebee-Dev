'use client';
import React from 'react';
import ApplicantInfo from '../../components/CreateZeeBee/ApplicantInfo';
import BorrowerInfo from '../../components/CreateZeeBee/BorrowerInfo';
import PropertyInfo from '../../components/CreateZeeBee/PropertyInfo';
import Loan from '../../components/CreateZeeBee/Loan';
import Experience from '../../components/CreateZeeBee/Experience';
import ThirdPartyService from '../../components/CreateZeeBee/ThirdPartyService';
import Declaration from '../../components/CreateZeeBee/Declaration';
import Demographics from '../../components/CreateZeeBee/Demographics';
import ReviewAndSign from '../../components/CreateZeeBee/ReviewAndSign';

const CommonZeeBeeForm = ({ form, pageKey }) => {
	const getContent = () => {
		switch (pageKey) {
			case 'borrowerInformation':
				return (
					<>
						<ApplicantInfo form={form} />
						<BorrowerInfo form={form} />
					</>
				);
			case 'propertyAndLoan':
				return (
					<>
						<PropertyInfo form={form} />
						<Loan form={form} />
					</>
				);
			case 'experienceAndThirdService':
				return (
					<>
						<Experience form={form} />
						<ThirdPartyService form={form} />
					</>
				);
			case 'declarationAndDemographics':
				return (
					<>
						<Declaration form={form} />
						<Demographics form={form} />
					</>
				);
			case 'reviewAndSign':
				return <ReviewAndSign form={form} />;
			default:
				return '';
		}
	};
	return getContent();
};
export default CommonZeeBeeForm;
