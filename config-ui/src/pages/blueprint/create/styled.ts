/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

import styled from 'styled-components';

export const Container = styled.div``;

export const Content = styled.div`
  margin-top: 36px;
  margin-bottom: 24px;
  font-size: 12px;

  .card + .card {
    margin-top: 24px;
  }

  h2 {
    margin: 0;
    padding: 0;
    font-size: 16px;
    font-weight: 600;
  }

  h3 {
    margin: 0 0 8px;
    padding: 0;
    font-size: 14px;
    font-weight: 600;
  }

  .back {
    display: flex;
    align-items: center;
    margin-bottom: 12px;
    color: #7497f7;
    cursor: pointer;

    span.bp4-icon {
      margin-right: 4px;
      cursor: pointer;
    }
  }

  .connection {
    display: flex;
    align-items: center;
    margin-bottom: 12px;

    span {
      margin-left: 8px;
      font-size: 14px;
      color: #292b3f;
      font-weight: 600;
    }
  }
`;
