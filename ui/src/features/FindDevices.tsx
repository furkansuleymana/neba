import { useState } from "react";
import { Button, Table, Checkbox } from "@mantine/core";

const devices = [
  {
    serial_number: "A12B34C56D78",
    model: "M1075-L",
    ip_address: "192.168.1.1",
    os_version: "11.2.68",
  },
  {
    serial_number: "A12B34C56D77",
    model: "M1075-L",
    ip_address: "192.168.1.1",
    os_version: "11.2.68",
  },
  {
    serial_number: "A12B34C56D76",
    model: "M1075-L",
    ip_address: "192.168.1.1",
    os_version: "11.2.68",
  },
  {
    serial_number: "A12B34C56D75",
    model: "M1075-L",
    ip_address: "192.168.1.1",
    os_version: "11.2.68",
  },
  {
    serial_number: "A12B34C56D74",
    model: "M1075-L",
    ip_address: "192.168.1.1",
    os_version: "11.2.68",
  },
];

export function FindDevices() {
  const [selectedRows, setSelectedRows] = useState<string[]>([]);

  const allSelected = selectedRows.length === devices.length;
  const someSelected = selectedRows.length > 0 && !allSelected;

  const toggleAll = () => {
    setSelectedRows(allSelected ? [] : devices.map((d) => d.serial_number));
  };

  const rows = devices.map((device) => (
    <Table.Tr
      key={device.serial_number}
      bg={
        selectedRows.includes(device.serial_number)
          ? "var(--mantine-color-blue-light)"
          : undefined
      }
    >
      <Table.Td>
        <Checkbox
          aria-label="Select device"
          checked={selectedRows.includes(device.serial_number)}
          onChange={(event) =>
            setSelectedRows(
              event.currentTarget.checked
                ? [...selectedRows, device.serial_number]
                : selectedRows.filter((sn) => sn !== device.serial_number)
            )
          }
        />
      </Table.Td>
      <Table.Td>{device.serial_number}</Table.Td>
      <Table.Td>{device.model}</Table.Td>
      <Table.Td>{device.ip_address}</Table.Td>
      <Table.Td>{device.os_version}</Table.Td>
    </Table.Tr>
  ));

  return (
    <>
      <Table highlightOnHover stickyHeader>
        <Table.Thead>
          <Table.Tr>
            <Table.Th>
              <Checkbox
                aria-label="Select all"
                checked={allSelected}
                indeterminate={someSelected}
                onChange={toggleAll}
              />
            </Table.Th>
            <Table.Th>Serial Number</Table.Th>
            <Table.Th>Model</Table.Th>
            <Table.Th>IP Address</Table.Th>
            <Table.Th>AXIS OS Version</Table.Th>
          </Table.Tr>
        </Table.Thead>
        <Table.Tbody>{rows}</Table.Tbody>
      </Table>{" "}
      <Button
        onClick={() => console.log(selectedRows)}
        disabled={selectedRows.length === 0}
        mt="md"
      >
        Log Selected Devices
      </Button>
    </>
  );
}
