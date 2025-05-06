import React, { useState } from 'react';
import { Typography, Card, Button, Space } from '@arco-design/web-react';
import { Editor } from '@/components/Editor';
import { ModalButton, ModalButtonPreview } from '@/components/Button/Modal';
import { LinkButton, Link, LinkText } from '@/components/Button/Link';

function Example() {
  return (
    <Card style={{ minHeight: '80vh' }}>
      <Typography.Title heading={6}>
        This is a very basic and simple page
      </Typography.Title>
      <Typography.Text>You can add content here :)</Typography.Text>
      <Space direction='vertical'>
        <LinkButton href='dashboard/workplace' text='noRefresh'></LinkButton>
        <LinkButton href='dashboard/workplace' text='refresh' refresh={true}></LinkButton>
        <LinkButton href='dashboard/workplace' text='newTab' newTab={true}></LinkButton>
        <Link href='/dashboard/workplace'> <Button>xx</Button> </Link>
        <Link href='/dashboard/workplace'>xx</Link>
        <LinkText href='/dashboard/workplace' text='xx'></LinkText>
      </Space>
      <EditorSample></EditorSample>
      <ModalSample></ModalSample>
    </Card>
  );
}


function EditorSample() {
  const [value, setValue] = useState("123")
  return (
    <>
      <Editor value={value} setValue={setValue} extra={<ModalButtonPreview modalTextBody={value}></ModalButtonPreview>}></Editor >
      <p>Current value is: </p>
      <pre>{value}</pre>
      <Space direction='vertical'>
        <Button type='primary' onClick={() => { setValue("xxxxx") }}>x</Button>
      </Space>
    </>
  )
}

function ModalSample() {
  return (
    <>
      <ModalButton btnBody="asdf"></ModalButton>

    </>
  )
}

export default Example;
