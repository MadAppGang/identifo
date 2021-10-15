import React from "react";
import { useSelector, useDispatch } from "react-redux";
import Button from "~/components/shared/Button";
import AddIcon from "~/components/icons/AddIcon";
import { fetchApplications } from "~/modules/applications/actions";
import ApplicationList from "./ApplicationList";
import ApplicationsPlaceholder from "./ApplicationsPlaceholder";
import useProgressBar from "~/hooks/useProgressBar";

const ApplicationsMainView = ({ history }) => {
  const dispatch = useDispatch();
  const { progress, setProgress } = useProgressBar();

  const applications = useSelector((s) => s.applications.list);

  const fetchData = async () => {
    setProgress(70);
    await dispatch(fetchApplications());
    setProgress(100);
  };

  React.useEffect(() => {
    fetchData();
  }, []);

  const initiateCreation = () => {
    history.push("/management/applications/new");
  };

  if (!applications.length && !progress) {
    return (
      <section className="iap-management-section">
        <ApplicationsPlaceholder onCreateApplicationClick={initiateCreation} />
      </section>
    );
  }

  return (
    <section className="iap-management-section">
      <p className="iap-management-section__title">
        Applications
        <Button Icon={AddIcon} onClick={initiateCreation}>
          Create application
        </Button>
      </p>

      <p className="iap-management-section__description">
        Setup an iOS, Android, Web, or Desktop application to use Identifo.
      </p>

      <ApplicationList loading={!!progress} applications={applications} />
    </section>
  );
};

export default ApplicationsMainView;
